package coinapi

import "time"

type Config struct {
	Start    string
	Stop     string
	CoinID   string
	Interval string
}

type TickerHistorical struct {
	TimePeriodStart time.Time `json:"time_period_start"`
	TimePeriodEnd   time.Time `json:"time_period_end"`
	TimeOpen        time.Time `json:"time_open"`
	TimeClose       time.Time `json:"time_close"`
	PriceOpen       float64   `json:"price_open"`
	PriceHigh       float64   `json:"price_high"`
	PriceLow        float64   `json:"price_low"`
	PriceClose      float64   `json:"price_close"`
	VolumeTraded    float64   `json:"volume_traded"`
	TradesCount     int       `json:"trades_count"`
}

type options struct {
	Start    time.Time `url:"time_start"`
	End      time.Time `url:"time_end"`
	Interval string    `url:"period_id,omitempty"`
}
