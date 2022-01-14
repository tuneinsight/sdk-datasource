package sdk

import (
	"github.com/google/uuid"
	"github.com/tuneinsight/sdk-datasource/pkg/models"
	"gorm.io/gorm"
)

// DataSourceBaseModel is the underlying generic model that all data sources inherit, it is the model saved in the database.
type DataSourceBaseModel struct {
	gorm.Model
	ID        models.DataSourceID `gorm:"primaryKey"`
	Name      string              `gorm:"uniqueIndex:udx_name;not null"`
	DeletedAt gorm.DeletedAt      `gorm:"uniqueIndex:udx_name"`
	Type      DataSourceType      `gorm:"not null"`
	Owner     string              `gorm:"not null"`
}

// NewDataSourceBaseModel creates a new DataSourceBaseModel instance given required fields.
func NewDataSourceBaseModel(id models.DataSourceID, owner, name string, dsType DataSourceType) *DataSourceBaseModel {
	dsm := new(DataSourceBaseModel)
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
func (ds *DataSourceBaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if ds.ID == "" {
		ds.ID = NewDataSourceID()
	}
	return
}

// GetName returns the name of the datasource
func (ds *DataSourceBaseModel) GetName() string {
	return ds.Name
}

// SetName sets the name of the datasource
func (ds *DataSourceBaseModel) SetName(name string) {
	ds.Name = name
}

// GetType returns the type of the data source
func (ds *DataSourceBaseModel) GetType() DataSourceType {
	return ds.Type
}

// SetType sets the type of the data source
func (ds *DataSourceBaseModel) SetType(t DataSourceType) {
	ds.Type = t
}

// GetID Return the id of the DataSource.
func (ds *DataSourceBaseModel) GetID() models.DataSourceID {
	return ds.ID
}

// SetID sets the ID of the DataSource
func (ds *DataSourceBaseModel) SetID(id models.DataSourceID) {
	ds.ID = id
}

// GetModel returns the underlying DataSource
func (ds *DataSourceBaseModel) GetModel() *DataSourceBaseModel {
	return ds
}

// SetModel sets the model of the data source
func (ds *DataSourceBaseModel) SetModel(dsm *DataSourceBaseModel) {
	*ds = *dsm
}

// GetOwner returns the owner of the data source
func (ds *DataSourceBaseModel) GetOwner() string {
	return ds.Owner
}

// NewDataSourceID generates a new data source id
func NewDataSourceID() models.DataSourceID {
	return models.DataSourceID(uuid.New().String())
}
