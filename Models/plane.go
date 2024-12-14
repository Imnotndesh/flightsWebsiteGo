package Models

type Plane struct {
	ID           int    `json:"pid"`
	RegNo        string `json:"reg_no"`
	HeadHostess  string `json:"h_hostess"`
	SubHostess   string `json:"s_hostess"`
	FirstClass   int    `json:"f_class"`
	EconomyClass int    `json:"e_class"`
	Capacity     int    `json:"capacity"`
	Pilot        string `json:"pilot"`
	Airline      string `json:"airline"`
}
