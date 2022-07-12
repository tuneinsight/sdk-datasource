package sdk

import (
	"github.com/sirupsen/logrus"
	"github.com/tuneinsight/sdk-datasource/pkg/models"
)

// DataSourceType defines a data source type, which uniquely identifies a data source plugin.
// A data source plugin must expose a variable of this type and with the same name (i.e., "DataSourceType") for GeCo to load it.
type DataSourceType string

// DataSourceFactory defines a TI Note data source factory, which is a function that can instantiate a DataSource.
// A data source plugin must expose a variable of this type and with the same name (i.e., "DataSourceFactory") for the TI Note to load it.
type DataSourceFactory func(dsc *DataSourceCore, config map[string]interface{}, dbManager *DBManager) (ds DataSource, err error)

// DataSource defines a TI Note data source, which is instantiated by a DataSourceFactory.
// All DataSource implementations should embed DataSourceCore.
type DataSource interface {
	SetID(id models.DataSourceID)
	GetID() models.DataSourceID

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

	// Query data source with a specific operation.
	// "jsonParameters" and "jsonResults" are both serialized JSON payloads.
	// "outputDataObjectsSharedIDs" maps output names of data object to their corresponding shared IDs.
	// "outputDataObjects" is a slice of data objects that were output by the query.
	Query(userID string, operation string, jsonParameters []byte, outputDataObjectsSharedIDs map[OutputDataObjectName]models.DataObjectSharedID) (jsonResults []byte, outputDataObjects []DataObject, err error)

	// Close is called to close all connections related to the DataSource or other instances.
	Close() error
}
