package cmd

import "github.com/spf13/cobra"

func addDBFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("db.connection-string", "", "Postgres connection string.")
	cmd.PersistentFlags().String("db.host", "localhost", "Database host")
	cmd.PersistentFlags().String("db.port", "5432", "Database port")
	cmd.PersistentFlags().String("db.sslmode", "disable", "Database sslmode")
	cmd.PersistentFlags().String("db.dbname", "simulator", "Database name")
	cmd.PersistentFlags().String("db.user", "core", "Database user")
	cmd.PersistentFlags().String("db.password", "password", "Database password")
	cmd.PersistentFlags().Bool("db.automigrate", true, "Auto run database migrations")
}

func addAPIFlags(cmd *cobra.Command) {
	cmd.Flags().String("api.port", "3001", "HTTP API port")
	cmd.Flags().Bool("api.dev-cors", false, "Enable development cors for HTTP API")
	cmd.Flags().String("api.dev-cors-host", "", "Allowed host for HTTP API dev cors")
}

func addPaprikaFlags(cmd *cobra.Command) {
	cmd.Flags().String("sow.paprika.coinid", "eth-ethereum", "CoinPaprika coin ID")
	cmd.Flags().String("sow.paprika.start", "1 jan 2020", "Date and time of scraping start")
	cmd.Flags().String("sow.paprika.stop", "31 dec 2020", "Date and time of scraping end, inclusive")
	cmd.Flags().String("sow.paprika.interval", "5m", "Price data interval (5m 10m 15m 30m 45m 1h 2h 3h 6h 12h 24h 1d 7d 14d 30d 90d 365d)")
}

func addSmartAlphaFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("log-senior-redeem", false, "Enable logs for SeniorRedeem function")
	cmd.Flags().Bool("log-senior-deposit", false, "Enable logs for SeniorDeposit function")
	cmd.Flags().Bool("senior-exit-enabled", false, "Enable Senior exits")
	cmd.Flags().Bool("junior-exit-enabled", false, "Enable Junior exits")
}
