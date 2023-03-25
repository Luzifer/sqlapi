package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		DSN            string `flag:"dsn" default:"" description:"MySQL DSN to connect to: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		fmt.Printf("mysqlapi %s\n", version)
		os.Exit(0)
	}

	router := mux.NewRouter()
	router.HandleFunc("/{database}", handleRequest).Methods(http.MethodPost)

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           router,
		ReadHeaderTimeout: time.Second,
	}

	logrus.WithFields(logrus.Fields{
		"addr":    cfg.Listen,
		"version": version,
	}).Info("mysqlapi started")

	if err = server.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("listening for HTTP")
	}
}
