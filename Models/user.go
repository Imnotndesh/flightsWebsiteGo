package Models

type User struct {
	Name         string `json:"Full_Name"`
	Username     string `json:"Username"`
	PasswordHash string `json:"Password_Hash"`
	Phone        string `json:"Phone"`
	Email        string `json:"Email"`
	ID           int    `json:"User_Id"`
}
type UpdateUserRequest struct {
	Username     string `json:"Username,omitempty"`
	PasswordHash string `json:"Password_Hash,omitempty"`
	Phone        string `json:"Phone,omitempty"`
	Email        string `json:"Email,omitempty"`
	Name         string `json:"Name,omitempty"`
}
