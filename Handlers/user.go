package Handlers

import (
	"database/sql"
	"net/http"
)

// UserDetailsHandler -> Filter requests based on methods
func UserDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		case http.MethodPost:
			getUserDetails(db, w, r)
		}
	}
}

// getUserDetails -> Returns user info from DB using JSON response
func getUserDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {

}

// UserRegistrationHandler -> Filter requests based on methods
func UserRegistrationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		case http.MethodPost:
			registerUser(db, w, r)
		}
	}
}

// registerUser -> adds new user to DB using passed JSON as user details source
func registerUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {

}

// UserEditHandler -> handles editing requests based on method
func UserEditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		case http.MethodPost:
			updateUserDetails(db, w, r)
		}
	}
}

// updateUserDetails -> Updates user information inside db with provided ones
func updateUserDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {

}
