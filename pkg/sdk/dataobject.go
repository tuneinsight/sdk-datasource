package sdk

import "github.com/tuneinsight/sdk-datasource/pkg/models"

// DataObject defines a data object to be produced by a DataSourcePlugin.
type DataObject struct {
	// OutputName is the output name of the DataObject.
	OutputName OutputDataObjectName

	// SharedID is the DataObjectSharedID that this DataObject should have.
	SharedID models.DataObjectSharedID

	// TODO: the DataObjet defines at the moment int or float values, vectors or matrices,
	// TODO: but it should be integrated with the definition and implementation of the datamanager's DataObject to be generic

	// the ones not being used should be left to nil
	IntValue  *int64
	IntVector []int64
	IntMatrix [][]int64

	FloatValue  *float64
	FloatVector []float64
	FloatMatrix [][]float64

	Columns []string
}

// OutputDataObjectName is a name for a data object that was output by a DataSourcePlugin.Query.
type OutputDataObjectName string
