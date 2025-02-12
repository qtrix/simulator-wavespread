package cmd

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	formatter "github.com/lacasian/logrus-module-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initLogging() {
	logging := viper.GetString("logging")

	if verbose {
		logging = "*=debug"
	}

	if vverbose {
		logging = "*=trace"
	}

	if logging == "" {
		logging = "*=info"
	}
	viper.Set("logging", logging)

	gin.SetMode(gin.DebugMode)

	modules := formatter.NewModulesMap(logging)
	if level, exists := modules["gin"]; exists {
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	} else {
		level := modules["*"]
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	}

	f, err := formatter.New(modules)
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(f)

	log.Debug("Debug mode")
}

func callPersistentPreRun(cmd *cobra.Command, args []string) {
	if parent := cmd.Parent(); parent != nil {
		if parent.PersistentPreRun != nil {
			parent.PersistentPreRun(parent, args)
		}
	}
}

func buildDBConnectionString() {
	if viper.GetString("db.connection-string") == "" {
		user := viper.GetString("db.user")
		pass := viper.GetString("db.password")

		p := fmt.Sprintf("host=%s port=%s sslmode=%s dbname=%s user=%s password=%s", viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.sslmode"), viper.GetString("db.dbname"), user, pass)
		viper.Set("db.connection-string", p)
	}
}
func mustGetSubconfig(v *viper.Viper, key string, out interface{}) {
	err := unmarshalSubconfig(v, key, out)
	if err != nil {
		log.Fatal(err)
	}
}

func unmarshalSubconfig(v *viper.Viper, key string, out interface{}) error {
	vc := subtree(v, key)
	if vc == nil {
		return errors.Errorf("key '%s' not found", key)
	}
	err := vc.Unmarshal(out)
	return err
}

func subtree(v *viper.Viper, name string) *viper.Viper {
	r := viper.New()
	for _, key := range v.AllKeys() {
		if strings.Index(key, name+".") == 0 {
			r.Set(key[len(name)+1:], v.Get(key))
		}
	}
	return r
}
