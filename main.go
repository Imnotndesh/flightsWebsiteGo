package main

import (
	"AirportAPI/Database"
	"AirportAPI/Handlers"
	"log"
	"net/http"
)

const (
	ApiPort = ":9080"
)

func main() {
	db, err := Database.InitDB()
	if err != nil {
		log.Panicln("Cannot access database", err)
	}

	mux := http.NewServeMux()
	// Routes definition and startup for API

	// User endpoints
	mux.HandleFunc("/user/me", Handlers.UserDetailsHandler(db))
	mux.HandleFunc("/user/me/edit", Handlers.UserEditHandler(db))
	mux.HandleFunc("/user/me/delete", Handlers.UserDeleteHandler(db))
	mux.HandleFunc("/user/register", Handlers.UserRegistrationHandler(db))
	mux.HandleFunc("/user/tickets", Handlers.UserTicketsHandler(db))
	mux.HandleFunc("/user/login", Handlers.AuthHandler(db))
	mux.HandleFunc("/user/me/top-up", Handlers.TopUpHandler(db))
	// Flights endpoints
	mux.HandleFunc("/flights", Handlers.FlightsInformationHandler(db))
	mux.HandleFunc("/flights/book", Handlers.TicketBookingHandler(db))
	//Admin endpoints
	mux.HandleFunc("/admins/login", Handlers.AdminLoginHandler(db))
	mux.HandleFunc("/admins/register", Handlers.AdminRegistrationHandler(db))
	// Admin view
	mux.HandleFunc("/admins/view/planes", Handlers.PlaneViewHandler(db))
	mux.HandleFunc("/admins/view/users", Handlers.UserViewHandler(db))
	mux.HandleFunc("/admins/view/flights", Handlers.FlightViewHandler(db))
	// Admin Edit
	mux.HandleFunc("/admins/edit/planes", Handlers.PlaneEditHandler(db))
	mux.HandleFunc("/admins/edit/users", Handlers.UserTableEditHandler(db))
	mux.HandleFunc("/admins/edit/flights", Handlers.FlightEditHandler(db))

	mux.HandleFunc("/admins/delete/plane", Handlers.PlaneDeletionHandler(db))
	mux.HandleFunc("/admins/delete/user", Handlers.MainUserDeleteHandler(db))
	mux.HandleFunc("/admins/delete/flight", Handlers.FlightDeletionHandler(db))

	log.Println("Starting server on port", ApiPort)
	err = http.ListenAndServe(ApiPort, mux)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
