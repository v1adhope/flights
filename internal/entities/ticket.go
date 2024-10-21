package entities

type Ticket struct {
	Id        string `json:"id" example:"uuid"`
	Provider  string `json:"provider" example:"Emirates"`
	FlyFrom   string `json:"flyFrom" example:"Moscow"`
	FlyTo     string `json:"flyTo" example:"Hanoi"`
	FlyAt     string `json:"flyAt" example:"3022-01-02T15:04:05+03:00"`
	ArriveAt  string `json:"arriveAt" example:"3022-01-03T18:04:40+07:00"`
	CreatedAt string `json:"createdAt" example:"timestampz"`
}

type TicketWholeInfo struct {
	Ticket
	Passengers []PassengerTicketWholeInfo `json:"passengers,omitempty"`
}
