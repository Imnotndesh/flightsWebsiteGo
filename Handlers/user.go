package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

type userRequest struct {
	username string
}

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
	var req userRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	var userDetails Models.User
	err := db.QueryRow("SELECT PHONE,EMAIL,FNAME,UID FROM users WHERE UNAME = ?", req.username).Scan(&userDetails.Phone, &userDetails.Email, &userDetails.Name, &userDetails.ID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userDetails)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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
	var newUserData Models.User
	if err := json.NewDecoder(r.Body).Decode(&newUserData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check for existing user
	var existingUserID int
	err := db.QueryRow("SELECT UID FROM users WHERE UNAME = ?", newUserData.Username).Scan(&existingUserID)
	if err == nil {
		http.Error(w, "User Exists", http.StatusInternalServerError)
		return
	}
	res, err := db.Exec("INSERT INTO users(uname, phone, email, fname, pass_hash) VALUES(?,?,?,?,?)", newUserData.Username, newUserData.Phone, newUserData.Email, newUserData.Name, newUserData.PasswordHash)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	userID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
	}
	newUserData.ID = int(userID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newUserData)
	if err != nil {
		http.Error(w, "Error returning new user details", http.StatusInternalServerError)
		return
	}

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
