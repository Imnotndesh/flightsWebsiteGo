package Models

type User struct {
	Name         string `json:"full_name,omitempty"`
	Username     string `json:"username,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	ID           int    `json:"user_id,omitempty"`
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
