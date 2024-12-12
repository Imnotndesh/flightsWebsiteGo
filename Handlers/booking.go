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

var (
	ticketInfo     Models.Ticket
	bookingRequest Models.BookingRequest
	youBrokeText   = "insufficient balance"
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
	var (
		bookingUser   Models.User
		usrNewBalance int
	)

	err = json.NewDecoder(r.Body).Decode(&bookingRequest)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Finding User information
	err = db.QueryRow("SELECT FNAME,UID,BALANCE FROM users WHERE UNAME = ? LIMIT 1", bookingRequest.Username).Scan(&ticketInfo.Name, ticketInfo.UID, &bookingUser.Balance)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// Finding flight information
	err = db.QueryRow("SELECT REGNO,DESTINATION,DEPATURE_TIME,AIRLINE,PRICE FROM flights WHERE FID = ?", bookingRequest.FlightID).Scan(&ticketInfo.RegNo, &ticketInfo.Destination, &ticketInfo.DepatureTime, &ticketInfo.Airline, &ticketInfo.Price)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	ticketInfo.FID, err = strconv.Atoi(bookingRequest.FlightID)

	// Check balance if actually available
	if bookingUser.Balance < ticketInfo.Price {
		http.Error(w, youBrokeText, http.StatusInternalServerError)
		return
	}
	usrNewBalance = bookingUser.Balance - ticketInfo.Price
	// Update the cut
	_, err = db.Exec("UPDATE users SET BALANCE = ? WHERE UNAME = ?", usrNewBalance, ticketInfo.Username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	ticketInfo.Price = ticketInfo.Price * bookingRequest.Tickets
	ticketInfo.Tickets = bookingRequest.Tickets
	// Storing gathered information into tickets table
	res, err := db.Exec("INSERT INTO tickets (REGNO, UID, FID, DEPATURE_TIME, FNAME, AIRLINE, DESTINATION,PRICE,TICKETS) VALUES (?,?,?,?,?,?,?,?,?)", ticketInfo.RegNo, ticketInfo.UID, ticketInfo.FID, ticketInfo.DepatureTime, ticketInfo.Name, ticketInfo.Airline, ticketInfo.Destination, ticketInfo.Price, ticketInfo.Tickets)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	num, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := res.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ticketInfo.ID = int(userID)
	if num == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Json response formulation
	w.WriteHeader(http.StatusOK)
}
