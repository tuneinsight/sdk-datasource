package sdk

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tuneinsight/sdk-datasource/pkg/models"
)

// DataSourceType defines a data source type, which uniquely identifies a data source plugin.
// A data source plugin must expose a variable of this type and with the same name (i.e., "DataSourceType") for GeCo to load it.
type DataSourceType string

// DataSourceFactory defines a TI Note data source factory, which is a function that can instantiate a DataSource.
// A data source plugin must expose a variable of this type and with the same name (i.e., "DataSourceFactory") for the TI Note to load it.
type DataSourceFactory func(dsc *DataSourceCore, config map[string]interface{}, dbManager *DBManager) (ds DataSource, err error)

const (
	// DefaultResultKey is the default key of the results returned by `Query`.
	DefaultResultKey string = "default"
)

// DataSource defines a TI Note data source, which is instantiated by a DataSourceFactory.
// All DataSource implementations should embed DataSourceCore.
type DataSource interface {
	SetID(id models.DataSourceID)
	GetID() models.DataSourceID

	SetContext(ctx *context.Context)
	GetContext() *context.Context

	// GetDataSourceCore returns the DataSourceCore of the DataSource.
	GetDataSourceCore() *DataSourceCore

	// SetDataSourceConfig sets the DataSource configuration, i.e. DataSource specific  metadata not contained in DataSourceCore.
	SetDataSourceConfig(map[string]interface{}) error
	// GetDataSourceConfig returns the DataSource configuration, i.e. DataSource specific  metadata not contained in DataSourceCore.
	GetDataSourceConfig() map[string]interface{}

	// Data must return all the DataSource data to be stored in the TI Note object storage.
	// Implementations should just call the DataImpl() function provided by the package.
	Data() map[string]interface{}

	// Config configures the DataSource. It must be called after DataSourceFactory.
	Config(logger logrus.FieldLogger, config map[string]interface{}) error

	// ConfigFromDB configures the DataSource with the info stored in the DB. It must be called after DataSourceFactory if the DataSource has been retrieved from the DB.
	ConfigFromDB(logger logrus.FieldLogger) error

	// Query data source with a specific operation. (the possible operations are defined by the data source)
	//  - userID is the ID of the user who is querying the data source (e.g. the username)
	//  - params are the query parameters. Typically two keys are "operation" for the type of operation
	// 	  and "params" for serialized JSON payloads.
	//  - resultKeys are the keys of the results to be returned. If empty, the default key is "default"
	Query(userID string, params map[string]interface{}, resultKeys ...string) (results map[string]interface{}, err error)

	// Close is called to close all connections related to the DataSource or other instances.
	Close() error
}
