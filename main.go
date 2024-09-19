package main

import (
	"AirportAPI/Database"
	"AirportAPI/Handlers"
	"log"
	"net/http"
)

const (
	ApiPort = "9080"
	WebPort = "9081"
)

func main() {
	db, err := Database.InitDB()
	if err != nil {
		log.Panicln("Cannot access database", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/user/me", Handlers.UserDetailsHandler(db))
	mux.HandleFunc("/user/me/edit", Handlers.UserEditHandler(db))
	mux.HandleFunc("/user/register", Handlers.UserRegistrationHandler(db))
	mux.HandleFunc("/user/tickets", Handlers.UserTicketsHandler(db))
	mux.HandleFunc("/flights", Handlers.FlightsInformationHandler(db))
	mux.HandleFunc("/flights/book", Handlers.TicketBookingHandler(db))
	err = http.ListenAndServe(":"+ApiPort, mux)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
