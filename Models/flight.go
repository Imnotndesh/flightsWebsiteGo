package Models

type Flight struct {
	ID             int    `json:"FID"`
	Destination    string `json:"DESTINATION"`
	Terminal       string `json:"TERMINAL"`
	Price          int    `json:"PRICE"`
	DepatureTime   int    `json:"DEPATURE_TIME"`
	Airline        string `json:"airline"`
	AvailableSeats int    `json:"AVAILABLE_SEATS"`
}

type Filters struct {
	MaxPrice    string `json:"max_price,omitempty"`
	MinPrice    string `json:"min_price,omitempty"`
	Destination string `json:"destination,omitempty"`
	Airline     string `json:"airline,omitempty"`
}
