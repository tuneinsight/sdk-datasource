package sdk

import "github.com/tuneinsight/sdk-datasource/pkg/models"

// DataObject defines a data object to be produced by a DataSourcePlugin.
type DataObject struct {
	// OutputName is the output name of the DataObject.
	OutputName OutputDataObjectName

	// SharedID is the DataObjectSharedID that this DataObject should have.
	SharedID models.DataObjectSharedID

	// todo: the DataObjet defines at the moment an int value and an int vector,
	// todo: but it should be integrated with the definition and implementation of the datamanager's DataObject to be generic

	// the one not being used should be left to nil
	IntValue  *int64
	IntVector []int64
}

// OutputDataObjectName is a name for a data object that was output by a DataSourcePlugin.Query.
type OutputDataObjectName string
