package Models

type Flight struct {
	ID             int    `json:"fid"`
	Destination    string `json:"destination"`
	Terminal       string `json:"terminal"`
	Price          int    `json:"price"`
	DepatureTime   string `json:"departure_time"`
	Airline        string `json:"airline"`
	AvailableSeats int    `json:"available_seats"`
	PID            int    `json:"pid"`
	REGNO          string `json:"regno"`
	Origin         string `json:"origin"`
}

type Filters struct {
	MaxPrice    string `json:"max_price,omitempty"`
	MinPrice    string `json:"min_price,omitempty"`
	Destination string `json:"destination,omitempty"`
	Airline     string `json:"airline,omitempty"`
}
