package coinapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/google/go-querystring/query"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/qtrix/simulator-wavespread/db"
	"github.com/qtrix/simulator-wavespread/utils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const (
	KeyLastDoneDay = "last-day-done"
)

var log = logrus.WithField("module", "tokenPrice")

type CoinApi struct {
	config     Config
	store      *db.DB
	start      time.Time
	stop       time.Time
	httpClient *http.Client
}

func New(config Config, store *db.DB, httpclient *http.Client) (*CoinApi, error) {
	var p CoinApi
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

	p.httpClient = httpclient
	p.store = store

	return &p, nil
}

func (p *CoinApi) insertPrices(ctx context.Context, tickers []*TickerHistorical) error {
	var rows [][]interface{}
	for _, t := range tickers {
		rows = append(rows, []interface{}{
			t.TimeClose, t.PriceClose,
		})
	}
	_, err := p.store.CopyFrom(ctx, pgx.Identifier{"price_data"}, []string{"timestamp", "price"}, pgx.CopyFromRows(rows))
	if err != nil {
		return errors.Wrap(err, "prices copy from")
	}

	return nil
}

func (p *CoinApi) buildURL(coinID string, options *options) (tickersHistorical []*TickerHistorical, err error) {
	baseURL := "https://rest.coinapi.io"
	//https://rest.coinapi.io/v1/ohlcv/BITSTAMP_SPOT_BTC_USD/history?period_id=5SEC&time_start=2023-01-10T13:05:07Z&time_end=2024-01-10T13:05:07Z
	url := fmt.Sprintf("%s/v1/ohlcv/%s/history", baseURL, coinID)
	url, err = constructURL(url, options)
	if err != nil {
		return nil, err
	}

	body, err := sendGET(p.httpClient, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &tickersHistorical)
	return tickersHistorical, err
}

func constructURL(rawURL string, options interface{}) (string, error) {
	if v := reflect.ValueOf(options); v.Kind() == reflect.Ptr && v.IsNil() {
		return rawURL, nil
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}

	values, err := query.Values(options)
	if err != nil {
		return rawURL, err
	}

	parsedURL.RawQuery = values.Encode()
	return parsedURL.String(), nil
}

func sendGET(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "text/json")
	req.Header.Add("X-CoinAPI-Key", "4AA09B40-629A-473D-A5B4-4307B62C57E0")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %v, body: %s", response.StatusCode, string(body))
	}

	return body, nil
}

func (p *CoinApi) Work(ctx context.Context) error {
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

		to := options{
			Start:    d,
			End:      d.AddDate(1, 0, 0),
			Interval: p.config.Interval,
		}

		// TODO check for 429 somehow
		tickers, err := p.buildURL(p.config.CoinID, &to)
		if err != nil {
			return errors.Wrap(err, "get historical tickers")
		}
		log.Println(tickers)

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
