package sdk

import (
	"time"

	"github.com/sirupsen/logrus"
)

// API is composed of a *sql.DB, logger and API configuration
type API struct {
	APIConfig
	logrus.FieldLogger
	MaxConnectionAttempts       int
	SleepingTimeBetweenAttempts time.Duration
}

// NewAPI creates a new API instance given configuration,connection, and parameters for connection attempts
func NewAPI(conf APIConfig) (*API, error) {
	api := new(API)
	api.APIConfig = conf
	api.FieldLogger = logrus.New().WithFields(logrus.Fields{})
	return api, nil
}
