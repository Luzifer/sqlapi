// Package types defines datatypes to work with
package types

import (
	"fmt"
)

type (
	// Query contains the query itself in the form the Adapter requires
	// and all arguments referenced in the query-string
	// (i.e. `[]any{"INSERT INTO foo VALUES (?, ?)", "bar", 1}`)
	Query []any

	// QueryResult contains the fields returned from the Query as a map.
	// The names of the fields are used as keys.
	QueryResult []map[string]any
)

// Args returns the arguments for the QueryString
func (q Query) Args() []any {
	if len(q) == 0 {
		return nil
	}
	return q[1:]
}

// QueryString returns the first argument as the query string
func (q Query) QueryString() (string, error) {
	if len(q) == 0 {
		return "", fmt.Errorf("no query given")
	}

	qs, ok := q[0].(string)
	if !ok {
		return "", fmt.Errorf("expected query as string in first argument, got %T", q[0])
	}

	return qs, nil
}
