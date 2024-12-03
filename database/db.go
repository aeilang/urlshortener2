package database

import (
	"database/sql"

	"github.com/aeilang/urlshortener/config"

	_ "github.com/lib/pq"
)

func NewDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.DSN()
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	// 设置数据库连接池
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
