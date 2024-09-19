package Handlers

import (
	"database/sql"
	"net/http"
)

func FlightsInformationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodPost:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		case http.MethodGet:
			getAvailableFlight(db, w, r)
		}
	}
}
func getAvailableFlight(db *sql.DB, w http.ResponseWriter, r *http.Request) {

}
