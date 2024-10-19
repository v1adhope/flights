package entities

type Ticket struct {
	Id        string `json:"id" example:"uuid"`
	Provider  string `json:"provider" example:"Emirates"`
	FlyFrom   string `json:"flyFrom" example:"Moscow"`
	FlyTo     string `json:"flyTo" example:"Hanoi"`
	FlyAt     string `json:"flyAt" example:"2022-01-02T12:04:05Z"`
	ArriveAt  string `json:"arriveAt" example:"2022-01-03T08:04:05Z"`
	CreatedAt string `json:"createdAt" example:"timestampz"`
}

type TicketWholeInfo struct {
	Ticket
	Passengers []PassengerTicketWholeInfo `json:"passengers,omitempty"`
}
