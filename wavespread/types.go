package wavespread

import "github.com/shopspring/decimal"

type AnchorPosition struct {
	Owner         string
	Amount        decimal.Decimal
	EntryPrice    decimal.Decimal
	EntryTime     int64
	ATokenBalance decimal.Decimal
	Exited        bool
}

type SurferPosition struct {
	Owner         string
	Amount        decimal.Decimal
	EntryTime     int64
	EntryPrice    decimal.Decimal
	STokenBalance decimal.Decimal
	Exited        bool
}
