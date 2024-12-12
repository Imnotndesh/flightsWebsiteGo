package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"net/http"
)

var (
	filters    Models.Filters
	flightInfo Models.Flight
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

// query: url:port/flights?q filter1=value ...
func getAvailableFlight(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var (
		flights []Models.Flight
		rows    *sql.Rows
	)
	rows, err = db.Query(`SELECT FID,DESTINATION,TERMINAL,PRICE,DEPATURE_TIME,AIRLINE,AVAILABLE_SEATS,REGNO,PID FROM flights`)
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		err = rows.Scan(&flightInfo.ID, &flightInfo.Destination, &flightInfo.Terminal, &flightInfo.Price, &flightInfo.DepatureTime, &flightInfo.Airline, &flightInfo.AvailableSeats, &flightInfo.REGNO, &flightInfo.PID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		flights = append(flights, flightInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flights)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
