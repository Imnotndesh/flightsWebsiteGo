package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// HashPassword Function to hash the password before storage
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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
func UserDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodDelete:
			deleteUser(db, w, r)
		}
	}
}

func deleteUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var (
		tbdUser Models.User
		err     error
	)

	if err = json.NewDecoder(r.Body).Decode(&tbdUser); err == nil || tbdUser.Username == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err = db.Exec(`DELETE FROM users WHERE UNAME = ?`, tbdUser.Username); err != nil || errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// getUserDetails -> Returns user info from DB using JSON response
func getUserDetails(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req Models.GetUserRequest
	)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	var fetchedUserDetails Models.User
	err = db.QueryRow("SELECT PHONE,EMAIL,FNAME,UID,BALANCE FROM users WHERE UNAME = ?", req.Username).Scan(&fetchedUserDetails.Phone, &fetchedUserDetails.Email, &fetchedUserDetails.Fullname, &fetchedUserDetails.ID, &fetchedUserDetails.Balance)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	fetchedUserDetails.Username = req.Username
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fetchedUserDetails)
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
func TopUpHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPut:
			TopUpBalance(db, w, r)
		}
	}
}
func IsUserInDB(username string, db *sql.DB) bool {
	var (
		userExists int
	)
	err := db.QueryRow(`SELECT users.UID FROM users WHERE UNAME = ?`, username).Scan(&userExists)
	if err != nil || errors.Is(err, sql.ErrNoRows) {
		return false
	}
	return true
}
func TopUpBalance(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var (
		topUpUser Models.User
	)

	if err := json.NewDecoder(r.Body).Decode(&topUpUser); err != nil || topUpUser.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !IsUserInDB(topUpUser.Username, db) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := db.Exec(`UPDATE users SET BALANCE = ? WHERE UNAME = ?`, topUpUser.Balance, topUpUser.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// registerUser -> adds new user to DB using passed JSON as user details source
func registerUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		newUserData Models.User
	)
	if err = json.NewDecoder(r.Body).Decode(&newUserData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check for existing user
	var existingUserID int
	err = db.QueryRow("SELECT UID FROM users WHERE UNAME = ?", newUserData.Username).Scan(&existingUserID)
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUserData.Balance = 0
	newUserData.PasswordHash, err = HashPassword(newUserData.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := db.Exec("INSERT INTO users(UNAME, PHONE, EMAIL, FNAME, PASS_HASH,BALANCE) VALUES(?,?,?,?,?,?)", newUserData.Username, newUserData.Phone, newUserData.Email, newUserData.Fullname, newUserData.PasswordHash, newUserData.Balance)
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
	var (
		editUserDetails Models.UpdateUserRequest
		newUserData     Models.User
	)
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
	hashedPass, err := HashPassword(editUserDetails.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUserData = Models.User{
		Username:     editUserDetails.Username,
		Fullname:     editUserDetails.Name,
		Phone:        editUserDetails.Phone,
		Email:        editUserDetails.Email,
		PasswordHash: hashedPass,
	}
	res, err := db.Exec("UPDATE users SET FNAME = ?, PHONE = ?, EMAIL = ? ,PASS_HASH = ? WHERE UNAME = ?", newUserData.Fullname, newUserData.Phone, newUserData.Email, newUserData.PasswordHash, newUserData.Username)
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

	// JSON response formulation
	w.WriteHeader(http.StatusOK)
}
