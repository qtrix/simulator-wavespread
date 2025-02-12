package smartalpha

import "github.com/shopspring/decimal"

type SeniorPosition struct {
	Owner         string
	Amount        decimal.Decimal
	EntryPrice    decimal.Decimal
	EntryTime     int64
	STokenBalance decimal.Decimal
	Exited        bool
}

type JuniorPosition struct {
	Owner         string
	Amount        decimal.Decimal
	EntryTime     int64
	EntryPrice    decimal.Decimal
	JTokenBalance decimal.Decimal
	Exited        bool
}
