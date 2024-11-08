package Models

type Ticket struct {
	ID           int    `json:"tid"`
	RegNo        string `json:"reg_no"`
	FID          int    `json:"fid"`
	UID          int    `json:"uid"`
	Username     string `json:"uname"`
	Destination  string `json:"destination"`
	DepatureTime int    `json:"departure_time"`
	Name         string `json:"fname"`
	Airline      string `json:"airline"`
	//backport this to ticket history frontend
	Price   int `json:"price"`
	Tickets int `json:"tickets"`
}
type BookingRequest struct {
	FlightID string `json:"flight_id"`
	Username string `json:"username"`
	Tickets  int    `json:"tickets"`
}
type UserTicketRequestFilters struct {
	UserId      string `json:"user_id"`
	Destination string `json:"destination"`
	Airline     string `json:"airline"`
}
