package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/qtrix/simulator-wavespread/db"
	"github.com/qtrix/simulator-wavespread/sower/coinapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var coinapiCmd = &cobra.Command{
	Use:   "coinapi",
	Short: "scrape coinapi for pricing data",
	Long:  "Address a wonderful greeting to the majestic executioner of this CLI",
	Run: func(cmd *cobra.Command, args []string) {
		// exit signals and cancellations
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		// config
		var coinapiCfg coinapi.Config
		mustGetSubconfig(viper.GetViper(), "sow.coinapi", &coinapiCfg)

		var dbCfg db.Config
		mustGetSubconfig(viper.GetViper(), "db", &dbCfg)

		// work
		store, err := db.New(ctx, dbCfg)
		if err != nil {
			log.Fatal(err)
		}

		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 3
		retryClient.HTTPClient.Timeout = time.Second * 10
		// retryClient.Logger = nil

		sower, err := coinapi.New(coinapiCfg, store, retryClient.HTTPClient)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			err := sower.Work(ctx)
			if err != nil {
				if ctx.Err() != nil {
					fmt.Printf("error:% v \n", ctx.Err())
				} else {
					fmt.Printf("error:% v \n", err)
				}
			}
			log.Info("sower worked finished")
		}()
		//spew.Dump(coinapiCfg)

		// exit
		<-ctx.Done()
		stop()
		log.Info("got stop signal, finishing work")

		//			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		//			defer cancel()
		// stop running stuff here
		// _ = ctx
		// select {
		// case :
		time.Sleep(time.Second * 5)
		log.Info("work done, goodbye")
		// case <-time.After(time.Second * 1):
		// 	fmt.Println("Timeout exiting")
		// }
		time.Sleep(time.Second * 5)
	},
}

func init() {
	sowCmd.AddCommand(coinapiCmd)

	addCoinApiFlags(coinapiCmd)
}
