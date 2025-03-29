package adapter

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	excludeDatabases = map[string]bool{
		"information_schema": true,
		"mysql":              true,
		"performance_schema": true,
		// "sys":                true,
	}
)

type MySQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type MySQLAdapter struct {
	status *Status
	config *MySQLConfig
	db     *sqlx.DB
}

func NewMySQLAdapter(config *MySQLConfig) (*MySQLAdapter, error) {
	return &MySQLAdapter{
		config: config,
		status: &Status{
			TableLoaded: map[string]bool{},
		},
	}, nil
}

func (a *MySQLAdapter) Status() *Status {
	return a.status
}

func (a *MySQLAdapter) Open() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", a.config.Username, a.config.Password, a.config.Host, a.config.Port, a.config.Database)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	a.db = db
	a.status.Opened = true
	return nil
}

func (a *MySQLAdapter) Close() error {
	err := a.db.Close()
	if err != nil {
		return err
	}
	a.status.Opened = false
	return nil
}

func (a *MySQLAdapter) ListDatabases() ([]string, error) {
	databases := []string{}
	err := a.db.Select(&databases, "SHOW DATABASES;")
	if err != nil {
		return nil, err
	}

	result := []string{}
	for i, db := range databases {
		if ok := excludeDatabases[db]; ok {
			continue
		}
		result = append(result, databases[i])
	}

	a.status.DatabaseLoaded = true

	return result, nil
}

func (a *MySQLAdapter) ListTables(database string) ([]string, error) {
	tables := []string{}
	err := a.db.Select(&tables, `
		SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = ?;
	`, database)
	if err != nil {
		return nil, err
	}
	a.status.TableLoaded[database] = true

	return tables, nil
}

// func (a *MySQLAdapter) ListTables() (map[string]TableInfo, error) {
// 	tables := []string{}
// 	err := a.db.Select(&tables, `
// 		SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES
// 		WHERE TABLE_SCHEMA = ?;
// 	`, a.config.Database)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result := make(map[string]TableInfo)

// 	for _, table := range tables {
// 		ti := TableInfo{Name: table}

// 		// Columns
// 		err = a.db.Select(&ti.Columns, `
// 			SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, EXTRA
// 			FROM INFORMATION_SCHEMA.COLUMNS
// 			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;
// 		`, a.config.Database, table)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Constraints
// 		err = a.db.Select(&ti.Constraints, `
// 			SELECT CONSTRAINT_NAME, CONSTRAINT_TYPE
// 			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
// 			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;
// 		`, a.config.Database, table)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Foreign keys
// 		err = a.db.Select(&ti.ForeignKeys, `
// 			SELECT CONSTRAINT_NAME, COLUMN_NAME, REFERENCED_TABLE_NAME, REFERENCED_COLUMN_NAME
// 			FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
// 			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND REFERENCED_TABLE_NAME IS NOT NULL;
// 		`, a.config.Database, table)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Indexes
// 		err = a.db.Select(&ti.Indexes, `
// 			SELECT INDEX_NAME, COLUMN_NAME, NON_UNIQUE
// 			FROM INFORMATION_SCHEMA.STATISTICS
// 			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?;
// 		`, a.config.Database, table)
// 		if err != nil {
// 			return nil, err
// 		}

// 		result[table] = ti
// 	}

// 	return result, nil
// }
