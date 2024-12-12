package Models

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fname"`
}
type DeleteRequest struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}
