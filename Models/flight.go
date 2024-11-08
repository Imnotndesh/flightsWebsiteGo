package Models

type Flight struct {
	ID             int    `json:"fid"`
	Destination    string `json:"destination"`
	Terminal       string `json:"terminal"`
	Price          int    `json:"price"`
	DepatureTime   int    `json:"departure_time"`
	Airline        string `json:"airline"`
	AvailableSeats int    `json:"available_seats"`
}

type Filters struct {
	MaxPrice    string `json:"max_price,omitempty"`
	MinPrice    string `json:"min_price,omitempty"`
	Destination string `json:"destination,omitempty"`
	Airline     string `json:"airline,omitempty"`
}
