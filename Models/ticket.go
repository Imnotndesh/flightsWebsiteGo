package Models

type Ticket struct {
	ID           int    `json:"tid"`
	RegNo        string `json:"reg_no"`
	FID          int    `json:"fid"`
	UID          int    `json:"uid"`
	Username     string `json:"uname"`
	Destination  string `json:"destination"`
	DepatureTime int    `json:"departure_time"`
	Name         string `json:"fmane"`
	Airline      string `json:"airline"`
}
type BookingRequest struct {
	FlightID string `json:"flight_id"`
	UserId   string `json:"user_id"`
}
