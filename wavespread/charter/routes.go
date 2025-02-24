package charter

import (
	"context"
	"fmt"
	"github.com/rs/cors"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.poolCompositionChart)
	mux.HandleFunc("/compare", c.comparisonChart)
	mux.HandleFunc("/surferToken", c.surfertokenPriceChart)
	mux.HandleFunc("/surfer-vs-anchor", c.surferVsAnchor)
	mux.HandleFunc("/profits-and-losses", c.profitsAndLosses)
	mux.HandleFunc("/liquidity", c.liquidity)

	// Enable CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://qtrix.github.io"}, // Replace with frontend URL
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	listenOn := fmt.Sprintf(":%s", c.APIConfig.Port)
	log.Infof("Listening on %s", listenOn)

	http.ListenAndServe(listenOn, handler)
}
