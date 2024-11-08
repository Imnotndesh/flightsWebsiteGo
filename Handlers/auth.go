package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	reqBody     Models.AuthRequest
	userDetails Models.User
	response    Models.ProcessingResponse
)

func AuthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			authenticateUser(db, w, r)
			return
		}
	}
}
func authenticateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = db.QueryRow("SELECT PASS_HASH FROM users WHERE UNAME=? LIMIT 1", reqBody.Username).Scan(&userDetails.PasswordHash)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(userDetails.PasswordHash), []byte(reqBody.Password)) != nil {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
