package sdk

import (
	"github.com/google/uuid"
	"github.com/tuneinsight/sdk-datasource/pkg/models"
	"gorm.io/gorm"
)

// DataSource is the underlying generic model that all data sources inherit, it is the model saved in the database.
type DataSourceModel struct {
	gorm.Model
	ID        models.DataSourceID `gorm:"primaryKey"`
	Name      string              `gorm:"uniqueIndex:udx_name;not null"`
	DeletedAt gorm.DeletedAt      `gorm:"uniqueIndex:udx_name"`
	Type      DataSourceType      `gorm:"not null"`
	Owner     string              `gorm:"not null"`
}

// NewDataSourceModel creates a new DataSourceModel instance given required fields.
func NewDataSourceModel(id models.DataSourceID, owner, name string, dsType DataSourceType) *DataSourceModel {
	dsm := new(DataSourceModel)
	if id == "" {
		dsm.ID = NewDataSourceID()
	} else {
		dsm.ID = id
	}
	dsm.Type = dsType
	dsm.Owner = owner
	dsm.Name = name
	return dsm
}

// BeforeCreate sets a new id to the data source if it is empty upon database insert.
func (ds *DataSourceModel) BeforeCreate(tx *gorm.DB) (err error) {
	if ds.ID == "" {
		ds.ID = NewDataSourceID()
	}
	return
}

// GetName returns the name of the datasource
func (ds *DataSourceModel) GetName() string {
	return ds.Name
}

// SetName sets the name of the datasource
func (ds *DataSourceModel) SetName(name string) {
	ds.Name = name
}

// GetType returns the type of the data source
func (ds *DataSourceModel) GetType() DataSourceType {
	return ds.Type
}

// SetType sets the type of the data source
func (ds *DataSourceModel) SetType(t DataSourceType) {
	ds.Type = t
}

// GetID Return the id of the DataSource.
func (ds *DataSourceModel) GetID() models.DataSourceID {
	return ds.ID
}

// SetID sets the ID of the DataSource
func (ds *DataSourceModel) SetID(id models.DataSourceID) {
	ds.ID = id
}

// GetModel returns the underlying DataSource
func (ds *DataSourceModel) GetModel() *DataSourceModel {
	return ds
}

// SetModel sets the model of the data source
func (ds *DataSourceModel) SetModel(dsm *DataSourceModel) {
	*ds = *dsm
}

// GetOwner returns the owner of the data source
func (ds *DataSourceModel) GetOwner() string {
	return ds.Owner
}

// NewDataSourceID generates a new data source id
func NewDataSourceID() models.DataSourceID {
	return models.DataSourceID(uuid.New().String())
}
