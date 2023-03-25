package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	request  [][]any
	response [][]map[string]any
)

/*
=== REQ
[
	["SELECT * FROM tablename WHERE name = ?", "foobar"]
]

=== RESP
[
	[{"name": "foobar", "age": 25}, {"name": "barfoo", "age": 56}]
]
*/

func executeQuery(db *sql.DB, query []any, resp *response) error {
	if len(query) == 0 {
		return errors.New("no query given")
	}

	qs, ok := query[0].(string)
	if !ok {
		return errors.Errorf("expected query as string in first argument, got %T", query[0])
	}

	rows, err := db.Query(qs, query[1:]...)
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	var respForQuery []map[string]any

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return errors.Wrap(err, "getting column types")
	}

	for rows.Next() {
		var (
			scanNames []string
			scanSet   []any
		)

		for _, col := range colTypes {
			scanNames = append(scanNames, col.Name())
			scanSet = append(scanSet, reflect.New(col.ScanType()).Interface())
		}

		if err = rows.Err(); err != nil {
			return errors.Wrap(err, "iterating rows")
		}

		if err = rows.Scan(scanSet...); err != nil {
			return errors.Wrap(err, "scanning row")
		}

		respForQuery = append(respForQuery, scanSetToObject(scanNames, scanSet))
	}

	if err = rows.Err(); err != nil {
		return errors.Wrap(err, "iterating rows (final)")
	}

	*resp = append(*resp, respForQuery)
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var (
		connID   = uuid.Must(uuid.NewV4()).String()
		database = mux.Vars(r)["database"]
		logger   = logrus.WithFields(logrus.Fields{
			"conn": connID,
			"db":   database,
		})

		connError = func(err error, reason string, code int) {
			logger.WithError(err).Error(reason)
			http.Error(w, fmt.Sprintf("an error occurred: %s", connID), http.StatusInternalServerError)
		}
	)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Conn-ID", connID)

	connInfo, err := mysql.ParseDSN(cfg.DSN)
	if err != nil {
		connError(err, "parsing DSN", http.StatusInternalServerError)
		return
	}
	connInfo.DBName = database

	db, err := sql.Open("mysql", connInfo.FormatDSN())
	if err != nil {
		connError(err, "opening db connection", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("closing db connection")
		}
	}()

	var (
		req  request
		resp response
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		connError(err, "parsing request", http.StatusBadRequest)
		return
	}

	for i, query := range req {
		if err = executeQuery(db, query, &resp); err != nil {
			connError(err, fmt.Sprintf("executing query %d", i), http.StatusInternalServerError)
			return
		}
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		connError(err, "encoding response", http.StatusInternalServerError)
		return
	}
}

//nolint:gocognit,gocyclo // contains simple type conversions
func scanSetToObject(scanNames []string, scanSet []any) map[string]any {
	row := make(map[string]any)
	for idx, name := range scanNames {
		// Some types are not very JSON friendly, lets make them
		switch tv := scanSet[idx].(type) {
		case *sql.NullBool:
			if tv.Valid {
				scanSet[idx] = tv.Bool
			} else {
				scanSet[idx] = nil
			}

		case *sql.NullByte:
			if tv.Valid {
				scanSet[idx] = tv.Byte
			} else {
				scanSet[idx] = nil
			}

		case *sql.NullFloat64:
			if tv.Valid {
				scanSet[idx] = tv.Float64
			} else {
				scanSet[idx] = nil
			}

		case *sql.NullInt16:
			if tv.Valid {
				scanSet[idx] = tv.Int16
			} else {
				scanSet[idx] = nil
			}

		case *sql.NullInt32:
			if tv.Valid {
				scanSet[idx] = tv.Int32
			} else {
				scanSet[idx] = nil
			}

		case *sql.NullInt64:
			if tv.Valid {
				scanSet[idx] = tv.Int64
			} else {
				scanSet[idx] = nil
			}

		case *sql.NullString:
			if tv.Valid {
				scanSet[idx] = tv.String
			} else {
				scanSet[idx] = nil
			}

		case *sql.NullTime:
			if tv.Valid {
				scanSet[idx] = tv.Time
			} else {
				scanSet[idx] = nil
			}

		case *sql.RawBytes:
			scanSet[idx] = string(*tv)
		}

		row[name] = scanSet[idx]
	}

	return row
}
