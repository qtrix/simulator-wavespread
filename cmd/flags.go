package cmd

import "github.com/spf13/cobra"

func addDBFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("db.connection-string", "postgres://uedecbulm85jjt:p6b28792047f99761090529f6628d7fea044ce2931930c966d6d5a86ffba5e637@c1i13pt05ja4ag.cluster-czrs8kj4isg7.us-east-1.rds.amazonaws.com:5432/da1jjt82naeovd", "Postgres connection string.")
	cmd.PersistentFlags().String("db.host", "localhost", "Database host")
	cmd.PersistentFlags().String("db.port", "5432", "Database port")
	cmd.PersistentFlags().String("db.sslmode", "disable", "Database sslmode")
	cmd.PersistentFlags().String("db.dbname", "simulator", "Database name")
	cmd.PersistentFlags().String("db.user", "core", "Database user")
	cmd.PersistentFlags().String("db.password", "password", "Database password")
	cmd.PersistentFlags().Bool("db.automigrate", false, "Auto run database migrations")
}

func addAPIFlags(cmd *cobra.Command) {
	cmd.Flags().String("api.port", "8080", "HTTP API port")
	cmd.Flags().Bool("api.dev-cors", true, "Enable development cors for HTTP API")
	cmd.Flags().String("api.dev-cors-host", "*", "Allowed host for HTTP API dev cors")
}

func addCoinApiFlags(cmd *cobra.Command) {
	cmd.Flags().String("sow.coinapi.coinid", "eth-ethereum", "coinapi coin ID")
	cmd.Flags().String("sow.coinapi.start", "1 jan 2024", "Date and time of scraping start")
	cmd.Flags().String("sow.coinapi.stop", "31 dec 2024", "Date and time of scraping end, inclusive")
	cmd.Flags().String("sow.coinapi.interval", "5m", "Price data interval (5m 10m 15m 30m 45m 1h 2h 3h 6h 12h 24h 1d 7d 14d 30d 90d 365d)")
}

func addWavespreadFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("log-anchor-redeem", false, "Enable logs for anchorRedeem function")
	cmd.Flags().Bool("log-anchor-deposit", false, "Enable logs for anchorDeposit function")
	cmd.Flags().Bool("anchor-exit-enabled", true, "Enable anchor exits")
	cmd.Flags().Bool("surfer-exit-enabled", true, "Enable surfer exits")
	cmd.Flags().Bool("surfer-switch-side-anchor-enabled", true, "Enable surfer switch side exits")
	cmd.Flags().Bool("anchor-switch-side-surfer-enabled", true, "Enable anchor switch side exits")
}
