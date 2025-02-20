package interfaces

import "github.com/shopspring/decimal"

type IOracle interface {
	PriceAtTs(ts int64) (decimal.Decimal, error)
}
