package main

import (
	"AirportAPI/Database"
	"AirportAPI/Handlers"
	"AirportAPI/Pages/Handler"
	"log"
	"net/http"
)

const (
	ApiPort = "9080"
	//WebPort = "9081"
)

func main() {
	db, err := Database.InitDB()
	if err != nil {
		log.Panicln("Cannot access database", err)
	}

	mux := http.NewServeMux()
	// Routes definition and startup for API
	mux.HandleFunc("/user/me", Handlers.UserDetailsHandler(db))
	mux.HandleFunc("/user/me/edit", Handlers.UserEditHandler(db))
	mux.HandleFunc("/user/register", Handlers.UserRegistrationHandler(db))
	mux.HandleFunc("/user/tickets", Handlers.UserTicketsHandler(db))
	mux.HandleFunc("/flights", Handlers.FlightsInformationHandler(db))
	mux.HandleFunc("/flights/book", Handlers.TicketBookingHandler(db))
	mux.HandleFunc("/user/login", Handlers.AuthHandler(db))

	// Pages Routes definition
	mux.HandleFunc("/login", Handler.LoginHandler())
	mux.HandleFunc("/register", Handler.RegistrationHandler())
	log.Println("Starting server on port :", ApiPort)
	err = http.ListenAndServe(":"+ApiPort, mux)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
