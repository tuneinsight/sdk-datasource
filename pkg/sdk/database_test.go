package sdk_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tuneinsight/sdk-datasource/pkg/sdk"
)

const (
	createQuery = `CREATE TABLE IF NOT EXISTS
	patients(name text, age int,weight real, height real)`
	insertQuery                 = `INSERT INTO patients (name,age,weight,height) VALUES (?,?,?,?)`
	retrieveQuery               = "SELECT * FROM patients"
	retrieveToFloatMatrix       = "SELECT weight, height FROM patients"
	name                        = "Albert Einstein"
	weight                      = 76.8
	height                      = 180.3
	age                         = 58
	maxConnAttempts             = 1
	sleepingTimeBetweenAttempts = 2
)

func TestSQLite(t *testing.T) {
	manager := sdk.NewDBManager(sdk.DBManagerConfig{
		MaxConnectionAttempts:              maxConnAttempts,
		SleepingTimeBetweenAttemptsSeconds: sleepingTimeBetweenAttempts,
	})
	dbConf := sdk.SQLiteConfig{
		Directory: t.TempDir(),
		Database:  "test",
	}

	db, err := manager.NewDatabase(&dbConf)
	require.NoError(t, err)
	_, err = db.Exec(createQuery)
	require.NoError(t, err)
	_, err = db.Exec(insertQuery, name, age, weight, height)
	require.NoError(t, err)

	rows, err := db.Retrieve(retrieveQuery)
	require.NoError(t, err)

	// Verify that it returned exactly one row and expected values match
	var actName string
	var actAge int
	var actWeight float64
	var actHeight float64
	require.True(t, rows.Next())
	require.NoError(t, rows.Scan(&actName, &actAge, &actWeight, &actHeight))
	require.Equal(t, name, actName)
	require.Equal(t, age, actAge)
	require.Equal(t, weight, actWeight)
	require.Equal(t, height, actHeight)
}
