package Models

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
