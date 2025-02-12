package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/coinpaprika/coinpaprika-api-go-client/coinpaprika"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/qtrix/simulator-wavespread/api"
	"github.com/qtrix/simulator-wavespread/db"
	"github.com/qtrix/simulator-wavespread/smartalpha/charter"
	"github.com/qtrix/simulator-wavespread/sower/paprika"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  "run sample command",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		var dbCfg db.Config
		mustGetSubconfig(viper.GetViper(), "db", &dbCfg)

		var paprikaCfg paprika.Config
		mustGetSubconfig(viper.GetViper(), "sow.paprika", &paprikaCfg)

		var apiCfg api.Config
		mustGetSubconfig(viper.GetViper(), "api", &apiCfg)

		store, err := db.New(ctx, dbCfg)
		if err != nil {
			log.Fatal(err)
		}

		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 3
		retryClient.HTTPClient.Timeout = time.Second * 10
		// retryClient.Logger = nil

		paprikaClient := coinpaprika.NewClient(retryClient.StandardClient())

		sower, err := paprika.New(paprikaClient, paprikaCfg, store)
		if err != nil {
			log.Fatal(err)
		}

		err = sower.Work(ctx)
		if err != nil {
			if ctx.Err() != nil {
				fmt.Printf("error:% v \n", ctx.Err())
			} else {
				fmt.Printf("error:% v \n", err)
			}
		}
		log.Info("sower worked finished")

		c := charter.New(store, apiCfg)
		go func() {
			c.Run(ctx)
		}()

		// exit
		<-ctx.Done()
		stop()
		log.Info("got stop signal, finishing work")
		log.Info("work done, goodbye")
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	addDBFlags(runCmd)
	addPaprikaFlags(runCmd)
	addSmartAlphaFlags(runCmd)
}
