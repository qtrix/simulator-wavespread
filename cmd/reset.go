package cmd

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/qtrix/simulator-wavespread/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset persistent storage",
	Long:  "use this comamnd  to erase all data",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var dbCfg db.Config
		err := unmarshalSubconfig(viper.GetViper(), "db", &dbCfg)
		if err != nil {
			log.Fatal(err)
		}
		dbCfg.AutoMigrate = false
		spew.Dump(dbCfg)

		// work
		psql, err := db.New(ctx, dbCfg)
		if err != nil {
			log.Fatal(err)
		}

		err = psql.ResetDB(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(resetCmd)

	addDBFlags(resetCmd)
}
