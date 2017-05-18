package server

const (
	// EngineLocal is constant for setting a local filesystem engine
	EngineLocal = "local"
	// EngineS3 is constant for setting an AWS S3 engine
	EngineS3 = "s3"
	// EngineSwift is constant for setting a swiftstack engine
	EngineSwift = "swift"
)

// Settings holds the configuration data for objstore
type Settings struct {
	// aws configuration settings
	Aws struct {
		AccessKey string
		SecretKey string
		Region    string
		Bucket    string
	}
	// engine type
	Engine string
	// local engine configuration
	Local struct {
		Root string
	}
	// newrelic configuration
	NewRelic struct {
		Appname string
		License string
		Enabled bool
	}
	// binding port for objstore
	Port int
	// swift engine configuration
	Swift struct {
		User      string `yaml:"apiuser"`
		Key       string `yaml:"apikey"`
		Container string
		AuthURL   string `yaml:"authurl"`
	}
}
