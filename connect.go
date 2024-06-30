package main

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func connect(database string) (db *sql.DB, err error) {
	switch cfg.DBType {
	case "mysql", "mariadb":
		connInfo, err := mysql.ParseDSN(cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("parsing DSN: %w", err)
		}
		connInfo.DBName = database

		if db, err = sql.Open("mysql", connInfo.FormatDSN()); err != nil {
			return nil, fmt.Errorf("opening db connection: %w", err)
		}

		return db, nil

	case "postgres", "pg", "crdb":
		u, err := url.Parse(cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("parsing DSN URL: %w", err)
		}

		u.Path = fmt.Sprintf("/%s", database)

		if db, err = sql.Open("pgx", u.String()); err != nil {
			return nil, fmt.Errorf("opening db connection: %w", err)
		}

		return db, nil

	default:
		return nil, fmt.Errorf("unknown database type %q", cfg.DBType)
	}
}
