package sdk

import (
	"fmt"
	"hash/crc64"
	"sync"
	"time"
)

// DBManager manages live database connections
type DBManager struct {
	DBManagerConfig
	sync.Mutex
	// maps parameters checksum to database interface
	databases map[int64]*Database
}

// DBManagerConfig regroups parameters for connection to all databases
type DBManagerConfig struct {
	MaxConnectionAttempts              int `yaml:"db-max-conn-attempts" default:"3"`
	SleepingTimeBetweenAttemptsSeconds int `yaml:"db-sleeping-time" default:"4"`
}

// NewDBManager creates a new manager with an empty map of databases
func NewDBManager(config DBManagerConfig) *DBManager {
	m := new(DBManager)
	m.DBManagerConfig = config
	m.databases = make(map[int64]*Database)
	return m
}

// NewDatabase creates a new database given a generic configuration
func (m *DBManager) NewDatabase(config DatabaseConfig) (db *Database, err error) {
	conn, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("creating database connection: %w", err)
	}
	db, err = NewDatabase(config, conn, m.MaxConnectionAttempts, time.Duration(m.SleepingTimeBetweenAttemptsSeconds)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("creating database: %w", err)
	}
	m.addDB(Checksum(config), db)
	return db, nil
}

// GetDatabase return the database instance from the configuration, or creates a new one if not registered
func (m *DBManager) GetDatabase(config DatabaseConfig) (db *Database, err error) {
	var ok bool
	db, ok = m.getDB(Checksum(config))
	if !ok {
		return m.NewDatabase(config)
	}
	return db, nil
}

// CloseAll attempts to close all db connections, concatenates any errors (locks the manager)
func (m *DBManager) CloseAll() error {
	m.Lock()
	defer m.Unlock()
	errs := make([]error, 0)
	for cs, db := range m.databases {
		err := db.Close()
		if err != nil {
			errs = append(errs, err)
		} else {
			delete(m.databases, cs)
		}
	}
	return WrapErrors("closing databases:", errs)
}

// Close closes the db given parameters
func (m *DBManager) Close(config DatabaseConfig) error {
	db, ok := m.getDB(Checksum(config))
	if !ok {
		return nil
	}
	err := db.Close()
	if err != nil {
		return err
	}
	m.deleteDB(Checksum(config))
	return nil
}

func (m *DBManager) addDB(cs int64, db *Database) {
	m.Lock()
	defer m.Unlock()
	m.databases[cs] = db
}

func (m *DBManager) getDB(cs int64) (db *Database, ok bool) {
	m.Lock()
	defer m.Unlock()
	db, ok = m.databases[cs]
	return
}

func (m *DBManager) deleteDB(cs int64) {
	m.Lock()
	defer m.Unlock()
	delete(m.databases, cs)
}

// Checksum returns the checksum of a connection configuration which is created using the datasource name and the driver name
func Checksum(config DatabaseConfig) int64 {
	uniqueString := config.DriverName() + "." + config.DataSourceName()
	return int64(crc64.Checksum([]byte(uniqueString), crc64.MakeTable(crc64.ECMA)))
}

// WrapErrors wraps multiple errors into a single error message
func WrapErrors(msg string, errs []error) error {
	if len(errs) <= 0 {
		return nil
	}
	errStr := msg + ":\n"
	for i, err := range errs {
		errStr += fmt.Sprintf("* %v", err.Error())
		if i < len(errs) {
			errStr += "\n"
		}
	}
	return fmt.Errorf(errStr)
}
