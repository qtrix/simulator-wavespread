package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/qtrix/simulator-wavespread/api"
	"github.com/qtrix/simulator-wavespread/wavespread/charter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/qtrix/simulator-wavespread/db"
)

var chartCmd = &cobra.Command{
	Use:   "chart",
	Short: "Generate pool composition chart for Wavespread",
	Run: func(cmd *cobra.Command, args []string) {
		// exit signals and cancellations
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		var dbCfg db.Config
		mustGetSubconfig(viper.GetViper(), "db", &dbCfg)

		var apiCfg api.Config
		mustGetSubconfig(viper.GetViper(), "api", &apiCfg)

		// work
		store, err := db.New(ctx, dbCfg)
		if err != nil {
			log.Fatal(err)
		}

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
	RootCmd.AddCommand(chartCmd)

	addDBFlags(chartCmd)
	addWavespreadFlags(chartCmd)
}
