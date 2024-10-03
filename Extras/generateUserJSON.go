package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type userDetails struct {
	name           string
	username       string
	email          string
	password       string
	hashedPassword string
	phone          string
}

func generateUserJSON(username string, name string, password string, email string, phone string) userDetails {
	newUserDetails := userDetails{
		name:     name,
		username: username,
		email:    email,
		password: password,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}
	newUserDetails.hashedPassword = string(hashedPassword)
	return newUserDetails
}
func main() {
	res := generateUserJSON("bishop", "bishopking", "igotking", "bishop@suv.com", "031456789")
	fmt.Println(res)
}
