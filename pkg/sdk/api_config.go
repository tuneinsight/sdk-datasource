package sdk

// APIConfig is a general interface for specific api configs, from the specific config, the API name and URL
// can be retrieved as well as the datasource name to use the API.
type APIConfig interface {
	// APIName should return the name of the for example "misp" or "generic"
	APIName() string
	GetURL() string
}

// MISPAPIConfig is the configuration for the MISP HTTP API
type MISPAPIConfig struct {
	URL   string `yaml:"api-url" default:"localhost"`
	Token string `yaml:"api-token" default:""`
}

// APIName returns the name of the driver which is 'misp'
func (conf MISPAPIConfig) APIName() string {
	return "misp"
}

// GetURL returns the url in configuration
func (conf MISPAPIConfig) GetURL() string {
	return conf.URL
}

// GenericAPIConfig is the configuration for a generic HTTP API
type GenericAPIConfig struct {
	URL      string `yaml:"api-url" default:"localhost"`
	User     string `yaml:"api-user" default:""`
	Password string `yaml:"api-pwd" default:""`
}

// APIName returns the name of the driver which is 'generic'
func (conf GenericAPIConfig) APIName() string {
	return "generic"
}

// GetURL returns the url in configuration
func (conf GenericAPIConfig) GetURL() string {
	return conf.URL
}
