package server

import (
	"github.com/mshindle/objstore/ops"
	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
)

func storageBuilder() error {
	var (
		e   ops.Engine
		err error
	)
	// select engine
	switch config.Engine {
	case EngineLocal:
		e = ops.NewLocalFile(config.Local.Root)
	case EngineS3:
		e = ops.NewS3(config.Aws.Region, config.Aws.Bucket)
	case EngineSwift:
		e, err = ops.NewSwiftEngine(config.Swift.User, config.Swift.Key, config.Swift.AuthURL, config.Swift.Container)
	default:
		logrus.WithField("engine", config.Engine).Error("unknown engine type specified")
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
