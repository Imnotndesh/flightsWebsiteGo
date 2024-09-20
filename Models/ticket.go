package Models

type Ticket struct {
	ID           int    `json:"TID"`
	PlaneRegno   int    `json:"REG_NO"`
	Username     string `json:"UNAME"`
	DepatureTime int    `json:"DEPATURE_TIME"`
	Name         string `json:"FNAME"`
	Airline      string `json:"AIRLINE"`
}
