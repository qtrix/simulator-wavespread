package api

type Config struct {
	Port        string
	DevCors     bool   `mapstructure:"dev-cors"`
	DevCorsHost string `mapstructure:"dev-cors-host"`
}
