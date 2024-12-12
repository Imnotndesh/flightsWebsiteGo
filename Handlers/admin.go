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
	adminUser     Models.AdminUser
	existingAdmin int
	isUnameSubmit bool
	deleteRequest Models.DeleteRequest
)

// Credential Handlers

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			AdminLogin(w, r, db)
			return
		}
	}
}
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			AdminRegistration(w, r, db)
			return
		}
	}
}

// View Handlers

func PlaneViewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			ViewPlaneData(w, r, db)
			return
		}
	}
}
func UserViewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			ViewUserData(w, r, db)
			return
		}
	}
}
func FlightViewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			ViewFlightData(w, r, db)
			return
		}
	}
}

// Edit Handlers

func PlaneEditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			EditPlanes(w, r, db)
			return
		}
	}
}

func FlightEditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			EditFlights(w, r, db)
			return
		}
	}
}

// Deletion Handlers

func UsersDeletionHandlers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			DeleteUser(w, r, db)
			return
		}
	}
}
func FlightDeletionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			DeleteFlight(w, r, db)
			return
		}
	}
}
func PlaneDeletionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPost:
			DeletePlane(w, r, db)
			return
		}
	}
}

// Functions for Admin Credential handling

func AdminLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var storedUser Models.AdminUser
	if (json.NewDecoder(r.Body).Decode(&adminUser)) == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = db.QueryRow("SELECT PASS_HASH from admins WHERE UNAME = ? LIMIT 1", adminUser.Username).Scan(&storedUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(adminUser.Password)) != nil {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
func AdminRegistration(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if (json.NewDecoder(r.Body)).Decode(&adminUser) == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var userExists int
	err = db.QueryRow("SELECT UID FROM admins WHERE UNAME = ? LIMIT 1", adminUser.Username).Scan(&userExists)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	adminUser.Password, err = HashPassword(adminUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO admins(UNAME,FNAME,PASSWORD) VALUES (?,?,?)", adminUser.Username, adminUser.Fullname, adminUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Functions to view table data
func isUnameInRequest(r *http.Request) (bool, string) {
	err = json.NewDecoder(r.Body).Decode(&adminUser)
	if err != nil {
		return false, adminUser.Username
	}
	if len(adminUser.Username) == 0 {
		return false, adminUser.Username
	}
	return true, adminUser.Username
}
func isAdminInDb(username string, db *sql.DB) bool {
	err = db.QueryRow("SELECT UID FROM admins WHERE UNAME = ? LIMIT 1", username).Scan(&existingAdmin)
	if err != nil {
		return false
	} else if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	return true
}
func ViewPlaneData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var (
		rows      *sql.Rows
		planeInfo Models.Plane
	)

	isUnameSubmit, adminUser.Username = isUnameInRequest(r)
	if !isUnameSubmit {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !(isAdminInDb(adminUser.Username, db)) {
		w.WriteHeader(http.StatusUnauthorized)
	}
	rows, err = db.Query("SELECT PID,REGNO,H_HOSTESS,S_HOSTESS,F_CLASS,E_CLASS,CAPACITY,PILOT,AIRLINE FROM planes")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var planeData []Models.Plane
	for rows.Next() {
		err = rows.Scan(&planeInfo.ID, &planeInfo.RegNo, &planeInfo.HeadHostess, &planeInfo.SubHostess, &planeInfo.FirstClass, &planeInfo.EconomyClass, &planeInfo.Capacity, &planeInfo.Pilot, &planeInfo.Airline)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		planeData = append(planeData, planeInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(planeData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func ViewUserData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var (
		rows     *sql.Rows
		userInfo Models.User
		userData []Models.User
	)

	isUnameSubmit, adminUser.Username = isUnameInRequest(r)
	if !isUnameSubmit {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !(isAdminInDb(adminUser.Username, db)) {
		w.WriteHeader(http.StatusUnauthorized)
	}
	rows, err = db.Query("SELECT UID, UNAME, FNAME,PHONE,EMAIL,BALANCE FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		err = rows.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Fullname, &userInfo.Phone, &userInfo.Email, &userInfo.Balance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userData = append(userData, userInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func ViewFlightData(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var (
		rows       *sql.Rows
		flightData []Models.Flight
	)

	isUnameSubmit, adminUser.Username = isUnameInRequest(r)
	if !isUnameSubmit {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !(isAdminInDb(adminUser.Username, db)) {
		w.WriteHeader(http.StatusUnauthorized)
	}
	rows, err = db.Query(`SELECT FID,DESTINATION,TERMINAL,PRICE,DEPATURE_TIME,AIRLINE,AVAILABLE_SEATS,REGNO,PID FROM flights`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&flightInfo.ID, &flightInfo.Destination, &flightInfo.Terminal, &flightInfo.Price, &flightInfo.DepatureTime, &flightInfo.Airline, &flightInfo.AvailableSeats, &flightInfo.REGNO, &flightInfo.PID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flightData = append(flightData, flightInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flightData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Editing table functions

func EditPlanes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var (
		planeInfo Models.Plane
	)

	if err = json.NewDecoder(r.Body).Decode(&planeInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if planeInfo.ID == 0 {
		_, err = db.Exec("INSERT INTO planes (REGNO,H_HOSTESS,S_HOSTESS,F_CLASS,E_CLASS,CAPACITY,PILOT,AIRLINE) VALUES (?,?,?,?,?,?,?,?)", planeInfo.RegNo, planeInfo.HeadHostess, planeInfo.SubHostess, planeInfo.FirstClass, planeInfo.EconomyClass, planeInfo.Capacity, planeInfo.Pilot, planeInfo.Airline)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	_, err = db.Exec("UPDATE planes SET REGNO = ?, H_HOSTESS = ?,S_HOSTESS = ?,F_CLASS = ?,E_CLASS = ?,CAPACITY = ?,PILOT = ?,AIRLINE = ? WHERE PID = ?", planeInfo.RegNo, planeInfo.HeadHostess, planeInfo.SubHostess, planeInfo.FirstClass, planeInfo.EconomyClass, planeInfo.Capacity, planeInfo.Pilot, planeInfo.Airline, planeInfo.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func EditFlights(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err = json.NewDecoder(r.Body).Decode(&flightInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if flightInfo.ID == 0 {
		_, err = db.Exec("INSERT INTO flights (DESTINATION,TERMINAL,PRICE,DEPATURE_TIME,AIRLINE,AVAILABLE_SEATS,REGNO,PID) VALUES (?,?,?,?,?,?,?,?)", flightInfo.Destination, flightInfo.Terminal, flightInfo.Price, flightInfo.DepatureTime, flightInfo.Airline, flightInfo.AvailableSeats, flightInfo.REGNO, flightInfo.PID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	_, err = db.Exec("UPDATE flights SET DESTINATION = ?,TERMINAL = ?,PRICE = ? ,DEPATURE_TIME = ?,AIRLINE = ?,AVAILABLE_SEATS = ?,REGNO = ? ,PID = ?  WHERE FID = ?", flightInfo.Destination, flightInfo.Terminal, flightInfo.Price, flightInfo.DepatureTime, flightInfo.Airline, flightInfo.AvailableSeats, flightInfo.REGNO, flightInfo.PID, flightInfo.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePlanes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err = json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil || deleteRequest.ID == 0 || len(deleteRequest.Username) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !(isAdminInDb(deleteRequest.Username, db)) {
		w.WriteHeader(http.StatusUnauthorized)
	}
	_, err = db.Exec("DELETE FROM planes WHERE PID = ?", deleteRequest.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func DeleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err = json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil || deleteRequest.ID == 0 || len(deleteRequest.Username) == 0 || !isAdminInDb(deleteRequest.Username, db) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM users WHERE UID = ?", deleteRequest.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func DeleteFlight(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err = json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil || deleteRequest.ID == 0 || len(deleteRequest.Username) == 0 || !isAdminInDb(deleteRequest.Username, db) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM flights WHERE FID = ?", deleteRequest.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func DeletePlane(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err = json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil || deleteRequest.ID == 0 || len(deleteRequest.Username) == 0 || !isAdminInDb(deleteRequest.Username, db) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM planes WHERE PID = ?", deleteRequest.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
