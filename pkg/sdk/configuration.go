package sdk

import (
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

// DatabaseConfig is a general interface for specific-db configs, from the specific config, the driver name
// can be retrieved as well as the datasource name to connect to the db and a function to register the actual driver
type DatabaseConfig interface {
	// DriverName should return the name of the driver for example "postgres" or "sqlite3"
	DriverName() string
	// DataSourceName should return the connection string to the db: postgres example "host=localhost port=5432 user=test password=test dbname=test sslmode=disable"
	DataSourceName() string
	// Driver Should return the database driver
	Driver() driver.Driver
	// Name should return the name of the connected database
	Name() string
}

// SQLiteConfig is the configuration when using the sqlite driver
type SQLiteConfig struct {
	Database  string `yaml:"db-database" default:"test"`
	Directory string `yaml:"db-directory" default:"db/"`
}

// DriverName returns the name of the driver which is sqlite3
func (conf SQLiteConfig) DriverName() string {
	return "sqlite3"
}

// DataSourceName should return the connection string to the db: <directory>/<database>.db
func (conf SQLiteConfig) DataSourceName() string {
	return conf.Directory + "/" + conf.Database + ".db"
}

// Driver returns the appropriate database driver
func (conf SQLiteConfig) Driver() driver.Driver {
	return &sqlite3.SQLiteDriver{}
}

// Name returns the name of the connected database
func (conf SQLiteConfig) Name() string {
	return conf.Database
}

// PostgresConfig is the configuration when using the postgres driver
type PostgresConfig struct {
	Host     string `yaml:"db-host" default:"localhost"`
	Port     int    `yaml:"db-port" default:"5432"`
	Database string `yaml:"db-database" default:"test"`
	User     string `yaml:"db-user" default:"postgres"`
	Password string `yaml:"db-pwd" default:"password"`
}

// DriverName returns "postgres"
func (conf PostgresConfig) DriverName() string {
	return "postgres"
}

// DataSourceName returns the connection string to a postgres db: "host=localhost port=5432 user=test password=test dbname=test sslmode=disable"
func (conf PostgresConfig) DataSourceName() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Database)
}

// Driver Should returns the postgres database driver
func (conf PostgresConfig) Driver() driver.Driver {
	return &pq.Driver{}
}

// Name should return the name of the connected database
func (conf PostgresConfig) Name() string {
	return conf.Database
}

// MySQLConfig is the configuration when using the MySQL driver
type MySQLConfig struct {
	Host     string `yaml:"db-host" default:"localhost"`
	Port     int    `yaml:"db-port" default:"5432"`
	Database string `yaml:"db-database" default:"test"`
	User     string `yaml:"db-user" default:"user"`
	Password string `yaml:"db-pwd" default:"password"`
}

// DriverName returns "mysql"
func (conf MySQLConfig) DriverName() string {
	return "mysql"
}

// DataSourceName returns the connection string to a MySQL db
func (conf MySQLConfig) DataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
}

// Driver Should returns the MySQL database driver
func (conf MySQLConfig) Driver() driver.Driver {
	return mysql.MySQLDriver{}
}

// Name should return the name of the connected database
func (conf MySQLConfig) Name() string {
	return conf.Database
}
