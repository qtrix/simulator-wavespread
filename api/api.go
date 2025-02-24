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

	// Manually configure CORS to allow all origins
	a.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*", "https://qtrix.github.io"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local testing
		fmt.Println("PORT environment variable is not set. Using default:", port)
	}

	a.setRoutes()

	logrus.Infof("starting api on port %s", port)

	err := a.engine.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *API) Close() {
}
