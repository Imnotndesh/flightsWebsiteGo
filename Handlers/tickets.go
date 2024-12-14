package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func UserTicketsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodGet:
			getUserTicketHistory(db, w, r)
		}
	}
}
func getUserTicketHistory(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var (
		ticketInfo Models.Ticket
		user       Models.User
		err        error
		query      url.Values
		tickets    []Models.Ticket
	)
	query = r.URL.Query()
	var queryUsername = query.Get("username")
	if queryUsername == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = db.QueryRow(`SELECT UID from users where UNAME = ?`, queryUsername).Scan(&user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	rows, err := db.Query(`SELECT TID,REGNO,FID,DEPATURE_TIME,AIRLINE,DESTINATION FROM tickets WHERE UID =?`, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ticketInfo.ID, &ticketInfo.RegNo, &ticketInfo.FID, &ticketInfo.DepatureTime, &ticketInfo.Airline, &ticketInfo.Destination)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			continue
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
