package server

import (
	"net/http"

	"github.com/newrelic/go-agent"

	"fmt"
	"time"

	"github.com/gocraft/web"
	"github.com/sirupsen/logrus"
)

// StoreContext holds object store contextual request information
type StoreContext struct {
	key string
}

var router *web.Router

func routeBuilder() error {

	// build our router
	router = web.New(StoreContext{})
	router.Middleware(loggerMiddleware)
	router.Middleware(ParseKey)
	router.Middleware(web.ShowErrorsMiddleware)

	router.Get(wrapHandle(relic, "/:*", GetObject))
	router.Put(wrapHandle(relic, "/:*", PutObject))
	router.Get("/", RootHandler)
	router.Put("/", RootHandler)

	return nil
}

func wrapHandle(app newrelic.Application, pattern string, fn func(*StoreContext, web.ResponseWriter, *web.Request)) (string, func(*StoreContext, web.ResponseWriter, *web.Request)) {
	return pattern, func(c *StoreContext, rw web.ResponseWriter, req *web.Request) {
		txn := app.StartTransaction(pattern, rw, req.Request)
		defer txn.End()

		fn(c, rw, req)
	}
}

// loggerMiddleware is generic middleware that will log requests to Logger (by default, Stdout).
func loggerMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	startTime := time.Now()

	next(rw, req)

	d := time.Since(startTime).Nanoseconds()

	var durationUnits string
	var duration int64
	switch {
	case d > 2000000:
		durationUnits = "ms"
		duration = d / 1000000
	case d > 1000:
		durationUnits = "μs"
		duration = d / 1000
	default:
		durationUnits = "ns"
		duration = d
	}

	logrus.WithFields(logrus.Fields{
		"duration":   d,
		"elapsed":    fmt.Sprintf("%d %s", duration, durationUnits),
		"statusCode": rw.StatusCode(),
		"method":     req.Method,
	}).Info(req.URL.Path)
}

// ParseKey pulls the storage key out of the request and stores it into context
func ParseKey(c *StoreContext, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.key = req.PathParams["*"]
	next(rw, req)
}

// GetObject retrieves an object from the storage using the URI Path as key.
// Leading slashes are stripped out. Getting "/" will return a bad request.
func GetObject(c *StoreContext, rw web.ResponseWriter, req *web.Request) {
	logrus.WithField("key", c.key).Info("starting GetObject")
	rw.Header().Set("Content-Type", "application/octet-stream")
	err := objstore.Retrieve(c.key, rw)
	if err != nil {
		logrus.WithField("key", c.key).Error("unable to read key from storage")
		http.NotFound(rw, req.Request)
		return
	}
}

// PutObject stores an object using the URI Path as the key.
// Leading slashes are stripped out. Putting an object to "/" will
// return a bad request.
func PutObject(c *StoreContext, rw web.ResponseWriter, req *web.Request) {
	logrus.WithField("key", c.key).Info("starting PutObject")
	err := objstore.Store(c.key, req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusPreconditionFailed)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}

// RootHandler takes care of bare root requests
func RootHandler(rw web.ResponseWriter, req *web.Request) {
	http.Error(rw, "cannot use / as a key", http.StatusBadRequest)
}
