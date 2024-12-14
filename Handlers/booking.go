package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func TicketBookingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			ticketBooking(db, w, r)
		}
	}
}
func ticketBooking(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var (
		err             error
		bookingUserCart []Models.BookingRequest
		// newBookingRequest Models.BookingRequest
	)
	err = json.NewDecoder(r.Body).Decode(&bookingUserCart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, newBookingRequest := range bookingUserCart {
		// find get flight information
		var (
			flightInfo      Models.Flight
			bookingUserInfo Models.User
		)
		err := db.QueryRow(`SELECT DESTINATION,PRICE,DEPATURE_TIME,REGNO,AIRLINE FROM flights WHERE FID = ? `, newBookingRequest.FlightID).Scan(&flightInfo.Destination, &flightInfo.Price, &flightInfo.DepatureTime, &flightInfo.REGNO, &flightInfo.Airline)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//find user information
		err = db.QueryRow(`SELECT BALANCE, FNAME, UID FROM users WHERE UNAME = ?`, newBookingRequest.Username).Scan(&bookingUserInfo.Balance, &bookingUserInfo.Fullname, &bookingUserInfo.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// checking if balance is enough for transaction
		totalPrice := newBookingRequest.Tickets * flightInfo.Price
		if bookingUserInfo.Balance < totalPrice {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			// Updating balance for user
			newBalance := bookingUserInfo.Balance - totalPrice
			_, err = db.Exec(`UPDATE users SET BALANCE = ? WHERE UID = ?`, newBalance, bookingUserInfo.ID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Println("Balance update successful")
		}

		// Insert ticket into db
		_, err = db.Exec(`INSERT INTO tickets (REGNO, UID, FID, DEPATURE_TIME, FNAME, AIRLINE, DESTINATION, PRICE, TICKETS,AIRLINE) VALUES (?,?,?,?,?,?,?,?,?,?)`, flightInfo.REGNO, bookingUserInfo.ID, newBookingRequest.FlightID, flightInfo.DepatureTime, bookingUserInfo.Fullname, flightInfo.Airline, flightInfo.Destination, totalPrice, newBookingRequest.Tickets, flightInfo.Airline)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}

}
