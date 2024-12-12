package Models

type User struct {
	Fullname     string `json:"full_name,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	ID           int    `json:"user_id,omitempty"`
	Balance      int    `json:"balance,omitempty"`
}
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"full_name,omitempty"`
}

type GetUserRequest struct {
	Username string `json:"Username"`
}
