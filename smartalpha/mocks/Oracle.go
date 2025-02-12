package mocks

import (
	"github.com/shopspring/decimal"
)

type Oracle struct {
	p map[int64]decimal.Decimal
}

func (o Oracle) PriceAtTs(ts int64) (decimal.Decimal, error) {
	return o.p[ts], nil
}

func (o *Oracle) SetPriceAtTs(ts int64, price decimal.Decimal) {
	if o.p == nil {
		o.p = make(map[int64]decimal.Decimal)
	}

	o.p[ts] = price
}
