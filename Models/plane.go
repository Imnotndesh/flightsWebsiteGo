package Models

type Plane struct {
	ID           int    `json:"PID"`
	RegNo        int    `json:"PREGNO"`
	HeadHostess  string `json:"H_HOSTESS"`
	SubHostess   string `json:"S_HOSTESS"`
	FirstClass   int    `json:"F_CLASS"`
	EconomyClass int    `json:"E_CLASS"`
	Capacity     int    `json:"CAPACITY"`
}
