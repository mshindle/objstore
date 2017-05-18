package ops

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ncw/swift"
)

// SwiftEngine defines a SwiftStack backed object storage engine
type SwiftEngine struct {
	connection *swift.Connection
	container  string
}

// NewSwiftEngine creates a Swiftstack based storage engine
func NewSwiftEngine(apiuser string, apikey string, authURL string, container string) (*SwiftEngine, error) {
	checkEnvDefault(&apiuser, "SWIFT_API_USER")
	checkEnvDefault(&apikey, "SWIFT_API_KEY")
	checkEnvDefault(&authURL, "SWIFT_AUTH_URL")
	checkEnvDefault(&container, "SWIFT_CONTAINER")

	c := &swift.Connection{
		UserName: apiuser,
		ApiKey:   apikey,
		AuthUrl:  authURL,
	}

	err := c.Authenticate()
	if err != nil {
		return nil, err
	}

	e := &SwiftEngine{
		connection: c,
		container:  container,
	}
	return e, nil
}

func checkEnvDefault(param *string, envvar string) {
	if *param == "" {
		*param = os.Getenv(envvar)
	}
}

// WriteTo reads key from Swift and writes the bytes to w
func (e *SwiftEngine) WriteTo(key string, w io.Writer) error {
	_, err := e.connection.ObjectGet(e.container, key, w, true, nil)
	return err
}

// ReadFrom reads data from r and stores it under key
func (e *SwiftEngine) ReadFrom(key string, r io.Reader) error {
	logrus.WithFields(logrus.Fields{"container": e.container, "key": key}).Debug("SwiftEngine writing to storage...")
	_, err := e.connection.ObjectPut(e.container, key, r, true, "", "", nil)
	return err
}

// Delete removes the object
func (e *SwiftEngine) Delete(key string) error {
	return e.connection.ObjectDelete(e.container, key)
}

