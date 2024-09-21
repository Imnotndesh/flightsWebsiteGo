package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type userRequest struct {
	Username string `json:"username"`
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
	log.Println(req)
	err := db.QueryRow("SELECT PHONE,EMAIL,FNAME,UID FROM users WHERE UNAME = ?", req.Username).Scan(&userDetails.Phone, &userDetails.Email, &userDetails.Name, &userDetails.ID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	userDetails.Username = req.Username
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
	log.Println(newUserData)
	// Check for existing user
	var existingUserID int
	err := db.QueryRow("SELECT UID FROM users WHERE UNAME = ?", newUserData.Username).Scan(&existingUserID)
	if err == nil {
		http.Error(w, "User Exists", http.StatusInternalServerError)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}
	res, err := db.Exec("INSERT INTO users(UNAME, PHONE, EMAIL, FNAME, PASS_HASH) VALUES(?,?,?,?,?)", newUserData.Username, newUserData.Phone, newUserData.Email, newUserData.Name, newUserData.PasswordHash)
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
		case http.MethodPut:
			updateUserDetails(db, w, r)
		}
	}
}

// updateUserDetails -> Updates user information inside db with provided ones
func updateUserDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var editUserDetails Models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&editUserDetails); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE UNAME = ?)", editUserDetails.Username).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	res, err := db.Exec("UPDATE users SET FNAME = ?, PHONE = ?, EMAIL = ? ,PASS_HASH = ? WHERE UNAME = ?", editUserDetails.Name, editUserDetails.Phone, editUserDetails.Email, editUserDetails.PasswordHash, editUserDetails.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Error updating user", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(editUserDetails)
	if err != nil {
		http.Error(w, "Error returning new user details", http.StatusInternalServerError)
		return
	}
}
