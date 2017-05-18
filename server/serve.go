package server

import (
	"fmt"
	"net/http"
	"errors"
	"log"
)

type builder func() error

var builders = []builder{
	validateConfigBuilder,
	storageBuilder,
	routeBuilder,
}

func ListenAndServe(settings *Settings) {
	connect := fmt.Sprintf(":%d", settings.Port)
	http.ListenAndServe(connect, router)
}

func validateConfigBuilder() error {
	// check port is valid
	if config.Port <= 0 {
		return errors.New("invalid port specified")
	}
	return nil
}

func storageBuilder() error {
	var (
		e   ops.Engine
		err error
	)
	// select engine
	switch config.Engine {
	case EngineLocal:
		e = engine.NewLocalFile(config.Local.Root)
	case EngineS3:
		e = engine.NewS3(config.Aws.Region, config.Aws.Bucket)
	case EngineSwift:
		e, err = engine.NewSwiftEngine(config.Swift.User, config.Swift.Key, config.Swift.AuthURL, config.Swift.Container)
	default:
		log.WithField("engine", config.Engine).Fatal("unknown engine type specified")
		err = errors.New("unknown engine type specified")
	}
	if err != nil {
		return err
	}

	// configure newrelic
	cfg := newrelic.NewConfig(config.NewRelic.Appname, config.NewRelic.License)
	relic, err = newrelic.NewApplication(cfg)
	if err != nil {
		return err
	}

	objstore = ops.NewStorage(&ops.Config{
		Engine: e,
		App:    relic,
	})
	return nil
}