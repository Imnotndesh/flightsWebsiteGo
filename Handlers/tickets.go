package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

var (
	conditions    []string
	args          []interface{}
	err           error
	query         url.Values
	searchFilters Models.UserTicketRequestFilters
	tickets       []Models.Ticket
)

func UserTicketsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodPost:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		case http.MethodGet:
			getUserTicketHistory(db, w, r)
		}
	}
}
func getUserTicketHistory(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	baseQuery := `SELECT TID,REGNO,FID,DEPATURE_TIME,AIRLINE,DESTINATION FROM tickets`
	query = r.URL.Query()

	// fetch primary filter
	searchFilters.UserId = query.Get("user_id")
	if searchFilters.UserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	conditions = append(conditions, "UID = ?")
	args = append(args, searchFilters.UserId)
	// Fetch optional filters
	searchFilters.Destination = query.Get("destination")
	searchFilters.Airline = query.Get("airline")

	if searchFilters.Airline == "" {
		conditions = append(conditions, "AIRLINE = ?")
		args = append(args, searchFilters.Airline)
	}
	if searchFilters.Destination == "" {
		conditions = append(conditions, "DESTINATION = ?")
		args = append(args, searchFilters.Destination)
	}
	if len(conditions) > 0 {
		baseQuery = baseQuery + " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		err = rows.Scan(&ticketInfo.ID, &ticketInfo.RegNo, &ticketInfo.FID, &ticketInfo.DepatureTime, &ticketInfo.Airline, &ticketInfo.Destination)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		tickets = append(tickets, ticketInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tickets)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
