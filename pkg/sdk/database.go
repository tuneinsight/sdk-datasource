package sdk

import (
	"database/sql"
	"reflect"
	"time"

	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// DatabaseError is wraps a database-related error
type DatabaseError struct {
	Err error
}

// Error prints the error message related to database
func (r *DatabaseError) Error() string {
	return fmt.Sprintf("database error: %v", r.Err)
}

// Database is composed of a *sql.DB, logger and Database configuration
type Database struct {
	DatabaseConfig
	logrus.FieldLogger
	*sql.DB
	MaxConnectionAttempts       int
	SleepingTimeBetweenAttempts time.Duration
}

// NewConnection opens a new sql.DB connection given the configuration, if the driver is not yet registered it gets registered
func NewConnection(conf DatabaseConfig) (*sql.DB, error) {
	if !IsRegistered(conf.DriverName()) {
		sql.Register(conf.DriverName(), conf.Driver())
	}
	return sql.Open(conf.DriverName(), conf.DataSourceName())
}

// NewDatabase creates a new Database instance given configuration,connection, and parameters for connection attempts
func NewDatabase(conf DatabaseConfig, connection *sql.DB, maxConnAttempts int, sleepingTimeBetweenAttempts time.Duration) (*Database, error) {
	var err error
	db := new(Database)
	db.DB = connection
	db.SleepingTimeBetweenAttempts = sleepingTimeBetweenAttempts
	db.MaxConnectionAttempts = maxConnAttempts
	db.DatabaseConfig = conf
	db.FieldLogger = logrus.New().WithFields(logrus.Fields{
		"db-driver": conf.DriverName(),
		"db-name":   conf.Name(),
	})
	err = db.WaitReady()
	if err != nil {
		return nil, &DatabaseError{Err: err}
	}
	return db, nil
}

// Store stores records in the Database.
// sqlStatement must be in the form of "INSERT INTO xxx VALUES ($1, $2, ...) RETURNING id".
func (db *Database) Store(sqlStatement string, args ...interface{}) (id string, err error) {

	err = db.WaitReady()

	if err != nil {
		return
	}

	err = db.QueryRow(sqlStatement, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("recording record in db: %w", err)
	}
	return
}

// Update updates records in the Database.
// sqlStatement must be in the form of "UPDATE xxx SET xxx WHERE xxx RETURNING id".
func (db *Database) Update(sqlStatement string, args ...interface{}) (id string, err error) {
	err = db.WaitReady()
	if err != nil {
		return
	}
	err = db.QueryRow(sqlStatement, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("updating record(s) in db: %w", err)
	}
	db.Debugf("correctly updated record with ID: %v", id)
	return
}

// Retrieve retrieves records from the Database.
// sqlStatement must be in the form of "SELECT xxx FROM xxx [WHERE xxx]".
func (db *Database) Retrieve(sqlStatement string, args ...interface{}) (rows *sql.Rows, err error) {
	err = db.WaitReady()
	if err != nil {
		return
	}
	rows, err = db.Query(sqlStatement, args...)
	if err != nil {
		return nil, fmt.Errorf("retrieving record(s) from db: %w", err)
	}
	return
}

// Delete deletes records from the Database.
// sqlStatement must be in the form of "DELETE FROM xxx WHERE xxx RETURN id".
func (db *Database) Delete(sqlStatement string, args ...interface{}) (err error) {
	err = db.WaitReady()
	if err != nil {
		return
	}
	var id string
	err = db.QueryRow(sqlStatement, args...).Scan(&id)
	if err != nil {
		return fmt.Errorf("deleting record(s) from db: %w", err)
	}
	return
}

// WaitReady performs a health-check of the Database and returns an error if the Database is unreachable
func (db *Database) WaitReady() (err error) {
	i := 0
	for err = db.Ping(); err != nil && i < db.MaxConnectionAttempts; err = db.Ping() {
		// Return directly if we encounter a MySQL error
		if reflect.TypeOf(err) == reflect.TypeOf(&mysql.MySQLError{}) {
			return err
		}

		// Otherwise, keep retrying (could be e.g. *net.OpError)
		db.FieldLogger.Warn(fmt.Errorf("impossible to connect to DB: %v. trying again in: %v seconds, error: %w", db.Name(), db.SleepingTimeBetweenAttempts, err))
		time.Sleep(db.SleepingTimeBetweenAttempts)
		i++
	}
	if err != nil {
		return fmt.Errorf("unable to connect to %v: %w", db.Name(), err)
	}
	return
}

// Close closes the sql.DB connection
func (db *Database) Close() error {
	return db.DB.Close()
}

// GetDBConn returns the sql.DB connection
func (db *Database) GetDBConn() *sql.DB {
	return db.DB
}

// IsRegistered returns whether or not the given db-driver is already registered
func IsRegistered(driver string) bool {
	registeredDrivers := sql.Drivers()
	for _, d := range registeredDrivers {
		if d == driver {
			return true
		}
	}
	return false
}
