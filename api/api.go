package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/qtrix/simulator-wavespread/db"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.WithField("module", "api")

type API struct {
	config Config
	engine *gin.Engine
	db     *db.DB
}

func New(config Config, db *db.DB) *API {
	return &API{
		config: config,
		db:     db,
	}
}

func (a *API) Run() {
	a.engine = gin.Default()

	a.engine = gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT environment variable is not set")
	}
	a.engine.Use(cors.Default())
	a.setRoutes()

	logrus.Infof("starting api on port %s", a.config.Port)
	err := a.engine.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *API) Close() {
}
