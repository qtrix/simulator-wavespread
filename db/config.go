package db

type Config struct {
	ConnectionString string `mapstructure:"connection-string"`
	AutoMigrate      bool
}
