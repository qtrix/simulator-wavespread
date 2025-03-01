package charter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qtrix/simulator-wavespread/api"
	"github.com/qtrix/simulator-wavespread/db"
	"github.com/qtrix/simulator-wavespread/wavespread/ethprice"
	"github.com/qtrix/simulator-wavespread/wavespread/interfaces"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "sa-charter")

type Charter struct {
	Oracle    interfaces.IOracle
	APIConfig api.Config
}

func New(db *db.DB, cfg api.Config) *Charter {
	oracle := ethprice.New(db)

	return &Charter{
		Oracle:    oracle,
		APIConfig: cfg,
	}
}

func (c *Charter) Run(ctx context.Context) {
	http.HandleFunc("/", c.poolCompositionChart)
	http.HandleFunc("/compare", c.comparisonChart)
	http.HandleFunc("/surferToken", c.surfertokenPriceChart)
	http.HandleFunc("/surfer-vs-anchor", c.surferVsAnchor)
	http.HandleFunc("/profits-and-losses", c.profitsAndLosses)
	http.HandleFunc("/liquidity", c.liquidity)

	listenOn := fmt.Sprintf(":%s", c.APIConfig.Port)
	log.Infof("listening on %s", listenOn)
	http.ListenAndServe(listenOn, nil)
}
