package datastore

import (
	"database/sql"

	"lince/entities"

	_ "github.com/go-sql-driver/mysql"
)

const defaultDSN = "root:juliano@tcp(localhost:3306)/stockops?parseTime=true"

// SettingsRepository fornece a conexão com o banco por empresa.
type SettingsRepository struct {
	Connection func(company entities.CompanyDatabaseConfig) *sql.DB
}

func NewSettingsRepository(cfg entities.Config) (SettingsRepository, error) {
	dsn := cfg.DSN
	if dsn == "" {
		dsn = defaultDSN
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return SettingsRepository{}, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return SettingsRepository{}, err
	}

	conn := func(company entities.CompanyDatabaseConfig) *sql.DB {
		_ = company
		return db
	}

	return SettingsRepository{Connection: conn}, nil
}
