package Handlers

import (
	"database/sql"
	"net/http"
)

func TicketBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			ticketBooking(db, w, r)
		}
	}
}
func ticketBooking(db *sql.DB, w http.ResponseWriter, r *http.Request) {

}
