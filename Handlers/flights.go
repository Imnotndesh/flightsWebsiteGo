package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
	baseQuery := `SELECT FID,DESTINATION,TERMINAL,PRICE,DEPATURE_TIME,AIRLINE,AVAILABLE_SEATS FROM flights`
	var conditions []string
	var args []interface{}
	var err error
	query := r.URL.Query()
	filters.Airline = query.Get("airline")
	filters.Destination = query.Get("destination")
	filters.MinPrice = query.Get("min_price")
	filters.MaxPrice = query.Get("max_price")
	if filters.Airline != "" {
		conditions = append(conditions, "AIRLINE = ?")
		args = append(args, filters.Airline)
	}
	if filters.Destination != "" {
		conditions = append(conditions, "DESTINATION = ?")
		args = append(args, filters.Destination)
	}
	if filters.MaxPrice != "" {
		conditions = append(conditions, "PRICE >= ?")
		args = append(args, filters.MaxPrice)
	}
	if filters.MinPrice != "" {
		conditions = append(conditions, "PRICE <= ?")
		args = append(args, filters.MinPrice)
	}
	if len(conditions) > 0 {
		baseQuery = baseQuery + " WHERE " + strings.Join(conditions, " AND ")
	}
	log.Println(baseQuery, args)
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
	var flights []Models.Flight
	for rows.Next() {
		err := rows.Scan(&flightInfo.ID, &flightInfo.Destination, &flightInfo.Terminal, &flightInfo.Price, &flightInfo.DepatureTime, &flightInfo.Airline, &flightInfo.AvailableSeats)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		flights = append(flights, flightInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flights)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
