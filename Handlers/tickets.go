package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

var (
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
	var user Models.User
	query = r.URL.Query()
	var queryUsername string = query.Get("Username")
	if queryUsername == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Println(queryUsername)
	err = db.QueryRow(`SELECT UID from users where UNAME = ?`, queryUsername).Scan(&user.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	log.Println(user.ID)
	rows, err := db.Query(`SELECT TID,REGNO,FID,DEPATURE_TIME,AIRLINE,DESTINATION FROM tickets WHERE UID =?`, user.ID)
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
