package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
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
	err = db.QueryRow("SELECT PASS_HASH,UID FROM users WHERE UNAME=? LIMIT 1", reqBody.Username).Scan(&userDetails.PasswordHash, &userDetails.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if userDetails.PasswordHash != reqBody.Password {
		response.Message = "Invalid Password"
		response.Success = false
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	} else {
		response.Success = true
		response.Message = "Login Successful"
		http.SetCookie(w, &http.Cookie{
			Name:    "localuid",
			Value:   strconv.Itoa(userDetails.ID),
			Path:    "/",
			Expires: time.Now().Add(10 * time.Minute),
			MaxAge:  600,
		})
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
