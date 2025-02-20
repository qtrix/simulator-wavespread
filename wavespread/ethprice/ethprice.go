package ethprice

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/qtrix/simulator-wavespread/db"
)

type ETHPrice struct {
	db *db.DB

	cache map[int64]decimal.Decimal
}

func (e *ETHPrice) PriceAtTs(ts int64) (decimal.Decimal, error) {
	start := time.Now()
	defer func() {
		logrus.WithField("duration", time.Since(start)).Trace("done fetching price")
	}()

	val, ok := e.cache[ts]
	if ok {
		return val, nil
	}

	// NOTE selecting first available price to deal with big coinpaprika gaps
	sel := `
		select price
		from price_data pd
		where pd.timestamp <= to_timestamp($1+3600)
		order by timestamp desc
		limit 1;
	`
	ctx := context.Background()
	row := e.db.QueryRow(ctx, sel, ts)
	var out decimal.Decimal
	err := row.Scan(&out)
	if err == pgx.ErrNoRows {
		return out, errors.New(fmt.Sprintf("price not found for given timestamp (%d)", ts))
	} else if err != nil {
		return out, errors.Wrapf(err, "fatching price at %d", ts)
	}

	e.cache[ts] = out

	return out, nil
}

func New(db *db.DB) *ETHPrice {
	return &ETHPrice{
		db:    db,
		cache: make(map[int64]decimal.Decimal),
	}
}
