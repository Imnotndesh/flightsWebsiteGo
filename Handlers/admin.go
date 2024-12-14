package Handlers

import (
	"AirportAPI/Models"
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var (
	adminUser     Models.AdminUser
	existingAdmin int
	isUnameSubmit bool
	deleteRequest Models.DeleteRequest
)

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
		case http.MethodPut:
			EditPlanes(w, r, db)
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
		case http.MethodPut:
			EditFlights(w, r, db)
		}
	}
}

// Deletion Handlers

func MainUserDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			fallthrough
		case http.MethodGet:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodDelete:
			DeleteUser(w, r, db)
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
		case http.MethodDelete:
			DeleteFlight(w, r, db)
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
		case http.MethodDelete:
			DeletePlanes(w, r, db)
		}
	}
}

// Functions for Admin Credential handling

func AdminLoginHandler(db *sql.DB) http.HandlerFunc {
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
func AdminLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var (
		err        error
		storedUser Models.AdminUser
	)
	if err = json.NewDecoder(r.Body).Decode(&adminUser); err != nil || adminUser.Username == "" {
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

func AdminRegistrationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AdminRegistration(w, r, db)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
func AdminRegistration(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var (
		newAdminUser Models.AdminUser
		err          error
	)
	err = json.NewDecoder(r.Body).Decode(&newAdminUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var userExists int
	err = db.QueryRow(`SELECT admins.UID FROM admins WHERE UNAME = ?`, newAdminUser.Username).Scan(&userExists)
	if userExists > 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newAdminUser.Password, err = HashPassword(newAdminUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO admins(UNAME, FNAME, PASS_HASH) VALUES (?,?,?)", newAdminUser.Username, newAdminUser.Fullname, newAdminUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Functions to view table data
func isUnameInRequest(r *http.Request) (bool, string) {
	err := json.NewDecoder(r.Body).Decode(&adminUser)
	if err != nil {
		return false, adminUser.Username
	}
	if len(adminUser.Username) == 0 {
		return false, adminUser.Username
	}
	return true, adminUser.Username
}
func isAdminInDb(username string, db *sql.DB) bool {
	err := db.QueryRow("SELECT UID FROM admins WHERE UNAME = ? LIMIT 1", username).Scan(&existingAdmin)
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
		return
	}
	rows, err := db.Query("SELECT PID,REGNO,H_HOSTESS,S_HOSTESS,F_CLASS,E_CLASS,CAPACITY,PILOT,AIRLINE FROM planes")
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
		err      error
	)

	isUnameSubmit, adminUser.Username = isUnameInRequest(r)
	if !isUnameSubmit {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !(isAdminInDb(adminUser.Username, db)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
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
		flightInfo Models.Flight
		err        error
	)

	isUnameSubmit, adminUser.Username = isUnameInRequest(r)
	if !isUnameSubmit {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !(isAdminInDb(adminUser.Username, db)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
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

	err := json.NewDecoder(r.Body).Decode(&planeInfo)
	if err != nil {
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
	var (
		flightInfo Models.Flight
		planeInfo  Models.Plane
	)
	err := json.NewDecoder(r.Body).Decode(&flightInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = db.QueryRow(`SELECT AIRLINE,REGNO,CAPACITY FROM planes WHERE PID = ?`, flightInfo.PID).Scan(&planeInfo.Airline, &planeInfo.RegNo, &planeInfo.Capacity)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if flightInfo.ID == 0 {
		_, err = db.Exec("INSERT INTO flights (DESTINATION,TERMINAL,PRICE,DEPATURE_TIME,AIRLINE,AVAILABLE_SEATS,REGNO,PID,ORIGIN) VALUES (?,?,?,?,?,?,?,?,?)", flightInfo.Destination, flightInfo.Terminal, flightInfo.Price, flightInfo.DepatureTime, planeInfo.Airline, planeInfo.Capacity, planeInfo.RegNo, flightInfo.PID, flightInfo.Origin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	_, err = db.Exec("UPDATE flights SET DESTINATION = ?,TERMINAL = ?,PRICE = ? ,DEPATURE_TIME = ?,AIRLINE = ?,AVAILABLE_SEATS = ?,REGNO = ? ,PID = ?, ORIGIN = ?  WHERE FID = ?", flightInfo.Destination, flightInfo.Terminal, flightInfo.Price, flightInfo.DepatureTime, planeInfo.Airline, flightInfo.AvailableSeats, planeInfo.RegNo, flightInfo.PID, flightInfo.Origin, flightInfo.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePlanes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var err error
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
	var err error
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
	var err error
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
func UserTableEditHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		case http.MethodPut:
			UserTableEdit(w, r, db)
		}
	}
}
func UserTableEdit(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var (
		tbeUser Models.User
		err     error
	)

	err = json.NewDecoder(r.Body).Decode(&tbeUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if tbeUser.ID == 0 {
		if tbeUser.PasswordHash, err = HashPassword(tbeUser.Password); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err := db.Exec(`INSERT INTO users (UNAME, PHONE, EMAIL, FNAME, BALANCE, PASS_HASH) VALUES (?,?,?,?,?,?)`, tbeUser.Username, tbeUser.Phone, tbeUser.Email, tbeUser.Fullname, tbeUser.Balance, tbeUser.PasswordHash)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else if tbeUser.ID != 0 {
		log.Println(tbeUser.Balance)
		_, err := db.Exec(`UPDATE users SET UNAME = ?, PHONE = ? ,EMAIL = ?, FNAME = ?, BALANCE = ?,PASS_HASH = ? WHERE UID = ?`, tbeUser.Username, tbeUser.Phone, tbeUser.Email, tbeUser.Email, tbeUser.Fullname, tbeUser.Balance, tbeUser.PasswordHash, tbeUser.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
