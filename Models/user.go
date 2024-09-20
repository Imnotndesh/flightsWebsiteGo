package Models

type User struct {
	Name         string `json:"FNAME"`
	Username     string `json:"UNAME"`
	PasswordHash string `json:"PASS_HASH"`
	Phone        string `json:"PHONE"`
	Email        string `json:"EMAIL"`
	ID           int    `json:"UID"`
}
