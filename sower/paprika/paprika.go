package paprika

import (
	"context"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/coinpaprika/coinpaprika-api-go-client/coinpaprika"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/qtrix/simulator-wavespread/db"
	"github.com/qtrix/simulator-wavespread/utils"
	"github.com/sirupsen/logrus"
)

const (
	KeyLastDoneDay = "last-day-done"
)

var log = logrus.WithField("module", "paprika")

type Paprika struct {
	config Config
	client *coinpaprika.Client
	store  *db.DB
	start  time.Time
	stop   time.Time
}

func (p *Paprika) Work(ctx context.Context) error {
	// check db for latest scrape
	var lastDate time.Time
	var gts db.StoreTime

	// TODO ugly as crap, make less so
	err := p.store.Get(ctx, KeyLastDoneDay, &gts)
	if err == pgx.ErrNoRows {
		lastDate = p.start
	} else if err != nil {
		return errors.Wrap(err, "get last day done")
	} else {
		lastDate = time.Time(gts).AddDate(0, 0, 1)
	}

	for d := lastDate; d.After(p.stop) == false; d = d.AddDate(0, 0, 1) {
		if utils.ContextIsDone(ctx) {
			return nil
		}
		fmt.Println(d.Format("2006-01-02"))

		to := coinpaprika.TickersHistoricalOptions{
			Start:    d,
			End:      d.AddDate(0, 0, 1),
			Interval: p.config.Interval,
		}
		tickers, err := p.client.Tickers.GetHistoricalTickersByID(p.config.CoinID, &to)
		// TODO check for 429 somehow
		if err != nil {
			return errors.Wrap(err, "get historical tickers")
		}
		log.Println(*tickers[0].Price)

		err = p.insertPrices(ctx, tickers)
		if err != nil {
			return errors.Wrap(err, "insert prices")
		}

		err = p.store.Set(ctx, KeyLastDoneDay, db.StoreTime(d))
		if err != nil {
			return errors.Wrap(err, "store last done day")
		}
	}

	return nil
}

func (p *Paprika) insertPrices(ctx context.Context, tickers []*coinpaprika.TickerHistorical) error {
	var rows [][]interface{}
	for _, t := range tickers {
		rows = append(rows, []interface{}{
			*t.Timestamp, *t.Price,
		})
	}
	_, err := p.store.CopyFrom(ctx, pgx.Identifier{"price_data"}, []string{"timestamp", "price"}, pgx.CopyFromRows(rows))
	if err != nil {
		return errors.Wrap(err, "prices copy from")
	}

	return nil
}

func New(client *coinpaprika.Client, config Config, store *db.DB) (*Paprika, error) {
	var p Paprika
	var err error

	p.config = config
	p.start, err = dateparse.ParseStrict(config.Start)
	if err != nil {
		return nil, errors.Wrap(err, "start date")
	}

	p.stop, err = dateparse.ParseStrict(config.Stop)
	if err != nil {
		return nil, errors.Wrap(err, "stop date")
	}

	log.Infof("scraping coin paprika from %s to %s", p.start, p.stop)

	p.client = client
	p.store = store

	return &p, nil
}
