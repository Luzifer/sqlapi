package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Luzifer/sqlapi/pkg/query"
	"github.com/Luzifer/sqlapi/pkg/types"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	request  []types.Query
	response []types.QueryResult
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
			http.Error(w, fmt.Sprintf("an error occurred: %s", connID), code)
		}
	)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Conn-ID", connID)

	db, err := connect(database)
	if err != nil {
		connError(err, "connecting to server", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("closing db connection")
		}
	}()

	var (
		req  request
		res  types.QueryResult
		resp response
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		connError(err, "parsing request", http.StatusBadRequest)
		return
	}

	for i, qry := range req {
		if res, err = query.RunQuery(db, qry); err != nil {
			connError(err, fmt.Sprintf("executing query %d", i), http.StatusInternalServerError)
			return
		}
		resp = append(resp, res)
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		connError(err, "encoding response", http.StatusInternalServerError)
		return
	}
}
