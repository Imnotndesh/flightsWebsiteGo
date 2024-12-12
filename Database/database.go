package Database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func InitDB() (*sql.DB, error) {
	dbpath := "db/airport.db"
	dbdir := "db"
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		err = os.Mkdir(dbdir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to database")
	err = createSchema(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func createSchema(db *sql.DB) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
	    UID INTEGER PRIMARY KEY AUTOINCREMENT,
	    UNAME TEXT UNIQUE,
	    PHONE TEXT UNIQUE,
	    EMAIL TEXT UNIQUE,
	    FNAME TEXT,
	    BALANCE INTEGER,
	    PASS_HASH TEXT UNIQUE
	);
`

	planesTable := `
	CREATE TABLE IF NOT EXISTS planes (
	    PID INTEGER PRIMARY KEY AUTOINCREMENT,
	    REGNO TEXT,
	    H_HOSTESS TEXT,
	    S_HOSTESS TEXT,
	    F_CLASS INTEGER,
	    E_CLASS INTEGER,
	    CAPACITY INTEGER,
	    PILOT TEXT,
	    AIRLINE TEXT
	);
`
	ticketTable := `
	CREATE TABLE IF NOT EXISTS tickets (
			TID  INTEGER PRIMARY KEY AUTOINCREMENT,
			REGNO TEXT,
			UID INT,
			FID INTEGER,
			DEPATURE_TIME TEXT,
			FNAME TEXT,
			AIRLINE TEXT,
			DESTINATION TEXT,
			PRICE INTEGER,
	    	FOREIGN KEY (UID) REFERENCES users(UID),
	    	FOREIGN KEY (FID) REFERENCES flights(FID)
	);
`
	adminsTable := `
CREATE TABLE IF NOT EXISTS admins (
    UNAME TEXT UNIQUE,
    FNAME TEXT,
    PASS_HASH TEXT,
)`
	flightTable := `
CREATE TABLE IF NOT EXISTS flights (
    FID INTEGER PRIMARY KEY AUTOINCREMENT,
    DESTINATION TEXT,
    TERMINAL TEXT,
    DEPATURE_TIME TEXT,
    PRICE INTEGER,
    AVAILABLE_SEATS INTEGER,
    AIRLINE TEXT,
    REGNO TEXT,
    PID INTEGER,
    FOREIGN KEY (PID) REFERENCES planes(PID)
);
`
	_, err := db.Exec(userTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(planesTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(ticketTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(flightTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(adminsTable)
	if err != nil {
		return err
	}
	log.Println("Successfully created database schema")
	return nil
}
