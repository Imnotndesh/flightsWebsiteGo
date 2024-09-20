package Models

type Flight struct {
	ID             int    `json:"FID"`
	Destination    string `json:"DESTINATION"`
	Terminal       string `json:"TERMINAL"`
	Price          int    `json:"PRICE"`
	DepatureTime   int    `json:"DEPATURE_TIME"`
	AvailableSeats int    `json:"AVAILABLE_SEATS"`
}
