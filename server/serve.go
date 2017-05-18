package server

import (
	"fmt"
	"net/http"

	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	"github.com/mshindle/objstore/ops"
)

type builder func() error

var (
	builders = []builder{
		validateConfigBuilder,
		storageBuilder,
		routeBuilder,
	}
	relic newrelic.Application
	objstore *ops.Storage
)

func ListenAndServe(settings *Settings) error {
	// set our server config
	config = settings

	// build the server
	err := buildServer()
	if err != nil {
		return err
	}

	// create connection string
	connect := fmt.Sprintf(":%d", config.Port)
	http.ListenAndServe(connect, router)
	return nil
}

func buildServer() error {
	for _, b := range builders {
		err := b()
		if err != nil {
			logrus.WithField("error", err).Error("could not build application")
			return err
		}
	}
	return nil
}
