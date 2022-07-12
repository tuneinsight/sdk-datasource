package sdk

import (
	"github.com/google/uuid"
	"github.com/tuneinsight/sdk-datasource/pkg/models"
	"github.com/tuneinsight/sdk-datasource/pkg/sdk/credentials"
	"gorm.io/gorm"
)

// DataSourceCore is the base struct that all data source implementations must embed.
type DataSourceCore struct {
	*DataSourceDatabaseModel
	*DataSourceMetadata
}

// DataSourceDatabaseModel contains the data source information that are stored in the TI Note database.
type DataSourceDatabaseModel struct {
	gorm.Model
	ID        models.DataSourceID `gorm:"primaryKey"`
	Name      string              `gorm:"uniqueIndex:udx_name;not null"`
	DeletedAt gorm.DeletedAt      `gorm:"uniqueIndex:udx_name"`
	Type      DataSourceType      `gorm:"not null"`
	Owner     string              `gorm:"not null"`
}

// DataSourceMetadata contains data source metadata fields that are stored in the TI Note object storage.
type DataSourceMetadata struct {
	CredentialsProvider credentials.Provider
	Attributes          []string
	ConsentType         models.DataSourceConsentType
}

// NewDataSourceCore instantiates a DataSourceCore with default DataSourceMetadata, given the original DataSourceDatabaseModel.
func NewDataSourceCore(model *DataSourceDatabaseModel) *DataSourceCore {
	dsc := new(DataSourceCore)
	dsc.DataSourceDatabaseModel = new(DataSourceDatabaseModel)
	if model != nil {
		dsc.DataSourceDatabaseModel = model
	}
	dsc.DataSourceMetadata = NewDataSourceMetadata()
	return dsc
}

// NewDataSourceDatabaseModel instantiates a DataSourceDatabaseModel given the required fields.
func NewDataSourceDatabaseModel(id models.DataSourceID, owner, name string, dsType DataSourceType) *DataSourceDatabaseModel {
	dsm := new(DataSourceDatabaseModel)
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

// NewDataSourceMetadata instantiates a default DataSourceMetadata.
func NewDataSourceMetadata() *DataSourceMetadata {
	meta := new(DataSourceMetadata)
	meta.Attributes = make([]string, 0)
	meta.ConsentType = models.DataSourceConsentTypeUnknown
	return meta
}

const (
	dsDSMetadataField = "ds-ti-note-metadata"
)

// GetMetadata returns the data source metadata.
func (ds *DataSourceCore) GetMetadata() map[string]interface{} {
	return map[string]interface{}{dsDSMetadataField: ds.DataSourceMetadata}
}

// Data returns all the data source data that must be stored in the TI Note object storage.
func (ds *DataSourceCore) Data() map[string]interface{} {
	return ds.GetMetadata()
}

// BeforeCreate sets a new id to the data source if it is empty upon database insert.
func (ds *DataSourceDatabaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if ds.ID == "" {
		ds.ID = NewDataSourceID()
	}
	return
}

// GetName returns the name of the data source.
func (ds *DataSourceDatabaseModel) GetName() string {
	return ds.Name
}

// SetName sets the name of the data source.
func (ds *DataSourceDatabaseModel) SetName(name string) {
	ds.Name = name
}

// GetType returns the type of the data source.
func (ds *DataSourceDatabaseModel) GetType() DataSourceType {
	return ds.Type
}

// SetType sets the type of the data source.
func (ds *DataSourceDatabaseModel) SetType(t DataSourceType) {
	ds.Type = t
}

// GetID Return the id of the DataSource.
func (ds *DataSourceDatabaseModel) GetID() models.DataSourceID {
	return ds.ID
}

// SetID sets the ID of the DataSource.
func (ds *DataSourceDatabaseModel) SetID(id models.DataSourceID) {
	ds.ID = id
}

// GetDatabaseModel returns the underlying DataSource.
func (ds *DataSourceDatabaseModel) GetDatabaseModel() *DataSourceDatabaseModel {
	return ds
}

// SetModel sets the model of the data source.
func (ds *DataSourceDatabaseModel) SetModel(dsm *DataSourceDatabaseModel) {
	*ds = *dsm
}

// GetOwner returns the owner of the data source.
func (ds *DataSourceDatabaseModel) GetOwner() string {
	return ds.Owner
}

// NewDataSourceID generates a new data source ID.
func NewDataSourceID() models.DataSourceID {
	return models.DataSourceID(uuid.New().String())
}
