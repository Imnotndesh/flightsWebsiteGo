package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"net/http"
)

func FlightsInformationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodPost:
			getAvailableFlight(db, w, r)
		case http.MethodGet:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
func getAvailableFlight(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var filters Models.Filters
	var flightInfo Models.Flight
	err := json.NewDecoder(r.Body).Decode(&filters)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if filters.Airline == "" || filters.Destination == "" || filters.MaxPrice == 0 || filters.MinPrice == 0 {
		err := db.QueryRow("SELECT DESTINATION , FID , TERMINAL , PRICE , DEPATURE_TIME, AVAILABLE_SEATS FROM flights").Scan(&flightInfo.Destination, &flightInfo.ID, &flightInfo.Terminal, &flightInfo.Price, &flightInfo.DepatureTime, &flightInfo.AvailableSeats)
		if err != nil {
			http.Error(w, "Error fetching flights information", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(flightInfo)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

}
