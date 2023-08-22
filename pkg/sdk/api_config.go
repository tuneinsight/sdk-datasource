package sdk

// APIConfig is a general interface for specific api configs, from the specific config, the API name and URL
// can be retrieved as well as the datasource name to use the API.
type APIConfig interface {
	// APIName should return the name of the for example "misp" or "generic"
	APIName() string
	// GetURL should return URL of the API
	GetURL() string
	// GetCert should return the name of the certificate used for this datasource (if applicable, otherwise "")
	GetCert() string
}

// GenericAPIConfig is the configuration for a generic HTTP API
type GenericAPIConfig struct {
	URL   string `yaml:"api-url" json:"api-url" default:"localhost"`
	Token string `yaml:"api-token" json:"api-token" default:""`
	Cert  string `yaml:"cert" json:"cert" default:""`
}

// APIName returns the name of the API which is 'generic'
func (conf GenericAPIConfig) APIName() string {
	return "generic"
}

// GetURL returns the url in configuration
func (conf GenericAPIConfig) GetURL() string {
	return conf.URL
}

// GetCert returns the name of the certificate used for this datasource (if applicable)
func (conf GenericAPIConfig) GetCert() string {
	return conf.Cert
}
