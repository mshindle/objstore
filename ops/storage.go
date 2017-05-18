package ops

import (
	"log"
	"io"
	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
)

// DefaultCapacity sets the initial capacity for a buffer
const DefaultCapacity = 1024 * 1024 * 4

const txnRetrieve = "ops.retrieve"
const txnStore = "ops.store"

// Engine is a specific implementation of Storage
type Engine interface {
	WriteTo(string, io.Writer) error
	ReadFrom(string, io.Reader) error
	Delete(string) error
}

// Storage is an implementation independent interface to underlying ops engines
type Storage struct {
	engine   Engine
	newrelic newrelic.Application
}

// Config handles configuration of the ops proxy
type Config struct {
	Engine Engine
	App    newrelic.Application
}

// NewStorage creates a new ops instance implementing engine.
func NewStorage(cfg *Config) *Storage {
	if cfg.Engine == nil {
		log.Fatalln("Cannot create a ops proxy without an engine")
	}
	if cfg.App == nil {
		nrConfig := newrelic.NewConfig("widget", "")
		nrConfig.Enabled = false
		app, err := newrelic.NewApplication(nrConfig)
		if err != nil {
			logrus.WithField("error", err).Fatalln("could not create dummy new relic app")
		}
		cfg.App = app
	}
	return &Storage{engine: cfg.Engine, newrelic: cfg.App}
}

// Retrieve pulls the data from under key and puts the contents into data.
func (s *Storage) Retrieve(key string, data io.Writer) error {
	txn := s.newrelic.StartTransaction(txnRetrieve, nil, nil)
	defer txn.End()

	err := s.engine.WriteTo(key, data)
	if err != nil {
		txn.NoticeError(err)
		return err
	}
	return nil
}

// RetrieveBytes pulls the data from under key and returns it as a byte array
func (s *Storage) RetrieveBytes(key string) ([]byte, error) {
	wb := NewWriteBuffer(make([]byte, 0, DefaultCapacity))
	err := s.Retrieve(key, wb)
	if err != nil {
		return nil, err
	}
	return wb.Bytes(), nil
}

// Store reads the data from reader and persists it under the given key
func (s *Storage) Store(key string, data io.Reader) error {
	txn := s.newrelic.StartTransaction(txnStore, nil, nil)
	defer txn.End()

	err := s.engine.ReadFrom(key, data)
	if err != nil {
		txn.NoticeError(err)
		return err
	}
	return nil
}

// Delete removes key from ops
func (s *Storage) Delete(key string) error {
	return s.engine.Delete(key)
}

