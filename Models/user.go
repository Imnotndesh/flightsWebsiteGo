package Models

type User struct {
	Name         string `json:"full_name"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	ID           int    `json:"user_id"`
}
type UpdateUserRequest struct {
	Username     string `json:"Username,omitempty"`
	PasswordHash string `json:"Password_Hash,omitempty"`
	Phone        string `json:"Phone,omitempty"`
	Email        string `json:"Email,omitempty"`
	Name         string `json:"Name,omitempty"`
}

type GetUserRequest struct {
	Username string `json:"Username"`
}
