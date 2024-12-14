package Models

type AdminUser struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}
type DeleteRequest struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}
