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
	if err := db.Ping(); err != nil {
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
	flightTable := ``
	userTable := ``
	ticketTable := ``
	planesTable := ``
	_, err := db.Exec(flightTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(userTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(ticketTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(planesTable)
	if err != nil {
		return err
	}
	log.Println("Successfully created database schema")
	return nil
}
