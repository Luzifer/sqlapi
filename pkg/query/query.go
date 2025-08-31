// Package query converts data from the query into the response types
package query

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Luzifer/sqlapi/pkg/types"
)

// RunQuery takes a Query containing at least the query-string as
// the first parameter and optionally any argument referenced. It
// runs the query using the connection stored inside the Adapter.
// The result then is parsed into the QueryResult form using the
// field names as keys and values as typed values.
func RunQuery(db *sql.DB, q types.Query) (types.QueryResult, error) {
	qs, err := q.QueryString()
	if err != nil {
		return nil, fmt.Errorf("getting query-string: %w", err)
	}

	rows, err := db.Query(qs, q.Args()...)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	var respForQuery types.QueryResult

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("getting column types: %w", err)
	}

	for rows.Next() {
		var (
			scanNames []string
			scanSet   []any
		)

		for _, col := range colTypes {
			t := col.ScanType()
			if t == nil {
				t = reflect.TypeFor[any]()
			}

			scanNames = append(scanNames, col.Name())
			scanSet = append(scanSet, reflect.New(t).Interface())
		}

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("iterating rows: %w", err)
		}

		if err = rows.Scan(scanSet...); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		respForQuery = append(respForQuery, scanSetToObject(scanNames, scanSet))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows (final): %w", err)
	}

	return respForQuery, nil
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
