package wavespread

import (
	"time"

	"github.com/shopspring/decimal"
)

func TimeToDecimal(t time.Time) decimal.Decimal {
	return decimal.NewFromInt(t.Unix())
}

func (sa *Wavespread) GetCurrentPrice() (decimal.Decimal, error) {
	return sa.oracle.PriceAtTs(sa.Clock.Now().Unix())
}
