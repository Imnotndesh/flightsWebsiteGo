package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

var ()

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
	var (
		bookingRequests []Models.BookingRequest
		bookingUser     Models.User
		usrNewBalance   int
		ticketInfo      Models.Ticket
		youBrokeText    = "insufficient balance"
	)

	err = json.NewDecoder(r.Body).Decode(&bookingRequests)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	for _, bookingRequest := range bookingRequests {
		err = db.QueryRow("SELECT FNAME,UID,BALANCE FROM users WHERE UNAME = ? LIMIT 1", bookingRequest.Username).Scan(&ticketInfo.Name, ticketInfo.UID, &bookingUser.Balance)
		if err != nil {
			// Handle individual booking errors (optional)
			log.Printf("Error processing booking for user %s: %v\n", bookingRequest.Username, err)
			continue
		} else if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// Finding flight information
		err = db.QueryRow("SELECT REGNO,DESTINATION,DEPATURE_TIME,AIRLINE,PRICE FROM flights WHERE FID = ?", bookingRequest.FlightID).Scan(&ticketInfo.RegNo, &ticketInfo.Destination, &ticketInfo.DepatureTime, &ticketInfo.Airline, &ticketInfo.Price)
		if err != nil {
			// Handle individual booking errors (optional)
			log.Printf("Error finding flight for booking %v: %v\n", bookingRequest, err)
			continue
		} else if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		ticketInfo.FID, err = strconv.Atoi(bookingRequest.FlightID)

		// Check balance if actually available
		if bookingUser.Balance < ticketInfo.Price {
			// Handle individual booking errors (optional)
			log.Printf("Insufficient balance for user %s: %v\n", bookingRequest.Username, youBrokeText)
			continue
		}
		usrNewBalance = bookingUser.Balance - ticketInfo.Price

		// Update the user's balance
		_, err = db.Exec("UPDATE users SET BALANCE = ? WHERE UNAME = ?", usrNewBalance, ticketInfo.Username)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ticketInfo.Price = ticketInfo.Price * bookingRequest.Tickets
		ticketInfo.Tickets = bookingRequest.Tickets

		// Insert each booking into the tickets table
		res, err := db.Exec("INSERT INTO tickets (REGNO, UID, FID, DEPATURE_TIME, FNAME, AIRLINE, DESTINATION,PRICE,TICKETS) VALUES (?,?,?,?,?,?,?,?,?)", ticketInfo.RegNo, ticketInfo.UID, ticketInfo.FID, ticketInfo.DepatureTime, ticketInfo.Name, ticketInfo.Airline, ticketInfo.Destination, ticketInfo.Price, ticketInfo.Tickets)
		if err != nil {
			// Handle individual booking errors (optional)
			log.Printf("Error inserting booking for user %s: %v\n", bookingRequest.Username, err)
			continue
		}

		// Check for successful insertion
		num, err := res.RowsAffected()
		if err != nil || num == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
