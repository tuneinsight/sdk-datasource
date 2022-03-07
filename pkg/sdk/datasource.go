package sdk

import (
	"github.com/sirupsen/logrus"
	"github.com/tuneinsight/sdk-datasource/pkg/models"
)

// DataSourceType defines a data source type, which uniquely identifies a data source plugin.
// A data source plugin must expose a variable of this type and with the same name (i.e., "DataSourceType") for GeCo to load it.
type DataSourceType string

// DataSourceFactory defines a GeCo data source factory, which is a function that can instantiate a DataSource.
// A data source plugin must expose a variable of this type and with the same name (i.e., "DataSourceFactory") for GeCo to load it.
type DataSourceFactory func(id models.DataSourceID, owner, name string) (ds DataSource, err error)

// DataSource defines a GeCo data source, which is instantiated by a DataSourceFactory.
type DataSource interface {
	GetModel() *DataSourceModel
	SetID(id models.DataSourceID)
	GetID() models.DataSourceID
	GetOwner() string
	Data() map[string]interface{}
	// Config configures the data source. It must be called after DataSourceFactory.
	Config(logger logrus.FieldLogger, config map[string]interface{}) error
	// ConfigFromDB configures the data source with the info stored in the DB. It must be called after DataSourceFactory if the data source has been retrieved from the DB.
	ConfigFromDB(logger logrus.FieldLogger) error
	// Query data source with a specific operation.
	// "jsonParameters" and "jsonResults" are both serialized JSON payloads.
	// "outputDataObjectsSharedIDs" maps output names of data object to their corresponding shared IDs.
	// "outputDataObjects" is a slice of data objects that were output by the query.
	Query(userID string, operation string, jsonParameters []byte, outputDataObjectsSharedIDs map[OutputDataObjectName]models.DataObjectSharedID) (jsonResults []byte, outputDataObjects []DataObject, err error)
}
