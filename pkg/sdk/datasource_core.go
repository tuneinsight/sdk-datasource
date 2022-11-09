package sdk

import (
	"context"

	"github.com/google/uuid"
	"github.com/tuneinsight/sdk-datasource/pkg/models"
	"github.com/tuneinsight/sdk-datasource/pkg/sdk/credentials"
	"gorm.io/gorm"
)

// DataSourceCore contains the common DataSource metadata.
// All DataSource implementations must embed DataSourceCore.
type DataSourceCore struct {
	*MetadataDB
	*MetadataStorage
	Ctx *context.Context
}

// MetadataDB contains the common DataSource metadata that are stored in the TI Note database.
type MetadataDB struct {
	gorm.Model
	ID                      models.DataSourceID      `gorm:"primaryKey"`
	Name                    string                   `gorm:"uniqueIndex:udx_name;not null"`
	Type                    DataSourceType           `gorm:"not null"`
	CredentialsProviderType credentials.ProviderType `gorm:"not null"`
	Owner                   string                   `gorm:"not null"`
	DeletedAt               gorm.DeletedAt           `gorm:"uniqueIndex:udx_name"`
}

// MetadataStorage contains the common DataSource metadata that are stored in the TI Note object storage.
type MetadataStorage struct {
	CredentialsProvider credentials.Provider
	Attributes          []string
	ConsentType         models.DataSourceConsentType
}

// NewDataSourceCore instantiates a DataSourceCore with the provided @mdb and @mds.
// If either @mdb or @mds are nil, they are set to default values.
func NewDataSourceCore(mdb *MetadataDB, mds *MetadataStorage) *DataSourceCore {
	dsc := new(DataSourceCore)

	if mds != nil {
		dsc.MetadataStorage = mds
	} else {
		dsc.MetadataStorage = NewMetadataStorage(nil)
	}

	if mdb != nil {
		dsc.MetadataDB = mdb
		// override any previously existing credentials provider type with the right one
		dsc.MetadataDB.CredentialsProviderType = dsc.MetadataStorage.CredentialsProvider.Type()
	} else {
		dsc.MetadataDB = NewMetadataDB("", "", "", "", dsc.MetadataStorage.CredentialsProvider.Type())
	}

	return dsc
}

// NewMetadataDB instantiates a MetadataDB given the required fields.
//
// cpType is the credential provider type (e.g. localCredentialsProvider)
func NewMetadataDB(id models.DataSourceID, owner, name string, dsType DataSourceType, cpType credentials.ProviderType) *MetadataDB {
	mdb := new(MetadataDB)
	if id == "" {
		mdb.ID = newDataSourceID()
	} else {
		mdb.ID = id
	}
	mdb.Type = dsType
	mdb.Owner = owner
	mdb.Name = name
	mdb.CredentialsProviderType = cpType
	return mdb
}

// NewMetadataStorage instantiates a default MetadataStorage.
// If a nil @cp is passed, a default Local one is set.
func NewMetadataStorage(cp credentials.Provider) *MetadataStorage {
	ms := new(MetadataStorage)
	ms.Attributes = make([]string, 0)
	ms.ConsentType = models.DataSourceConsentTypeUnknown
	if cp != nil {
		ms.CredentialsProvider = cp
	} else {
		ms.CredentialsProvider = credentials.NewLocal(nil)
	}
	return ms
}

const (
	// DSCoreMetadataField is the key under which the MetadataStorage are stored in the TI Note storage.
	DSCoreMetadataField = "ds-core-metadata"
)

// GetDataSourceCore returns the DataSourceCore of the data source.
func (dsc *DataSourceCore) GetDataSourceCore() *DataSourceCore {
	return dsc
}

// BeforeCreate sets a new id to the DataSource if it is empty upon database insert.
func (mdb *MetadataDB) BeforeCreate(tx *gorm.DB) (err error) {
	if mdb.ID == "" {
		mdb.ID = newDataSourceID()
	}
	return
}

// GetID Return the id of the DataSource.
func (mdb *MetadataDB) GetID() models.DataSourceID {
	return mdb.ID
}

// SetID sets the ID of the DataSource.
func (mdb *MetadataDB) SetID(id models.DataSourceID) {
	mdb.ID = id
}

// Data returns MetadataStorage in the format expected by the TI Node to store it in storage.
func (ms *MetadataStorage) Data() map[string]interface{} {
	return map[string]interface{}{DSCoreMetadataField: ms}
}

// DataImpl is the function that should be called by all Data() implementations.
func DataImpl(ds DataSource) map[string]interface{} {
	data := map[string]interface{}{DSCoreMetadataField: ds.GetDataSourceCore().MetadataStorage}
	for k, v := range ds.GetDataSourceConfig() {
		data[k] = v
	}
	return data
}

// newDataSourceID generates a new DataSource ID.
func newDataSourceID() models.DataSourceID {
	return models.DataSourceID(uuid.New().String())
}
