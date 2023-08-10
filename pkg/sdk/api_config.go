package sdk

// APIConfig is a general interface for specific api configs, from the specific config, the API name and URL
// can be retrieved as well as the datasource name to use the API.
type APIConfig interface {
	// APIName should return the name of the for example "misp" or "generic"
	APIName() string
	// GetURL should return URL of the API
	GetURL() string
	// GetTeleportCert should return the name of the Teleport certificate used for this datasource (if applicable, otherwise "")
	GetTeleportCert() string
}

// GenericAPIConfig is the configuration for a generic HTTP API
type GenericAPIConfig struct {
	URL          string `yaml:"api-url" json:"api-url" default:"localhost"`
	Token        string `yaml:"api-token" json:"api-token" default:""`
	TeleportCert string `yaml:"teleport-cert" json:"teleport-cert" default:""`
}

// APIName returns the name of the API which is 'generic'
func (conf GenericAPIConfig) APIName() string {
	return "generic"
}

// GetURL returns the url in configuration
func (conf GenericAPIConfig) GetURL() string {
	return conf.URL
}

// GetTeleportCert returns the name of the Teleport certificate used for this datasource (if applicable)
func (conf GenericAPIConfig) GetTeleportCert() string {
	return conf.TeleportCert
}
