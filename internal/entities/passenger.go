package entities

type Passenger struct {
	Id         string `json:"id" example:"uuid"`
	FirstName  string `json:"firstName" example:"Wendi"`
	LastName   string `json:"lastName" example:"Reyes"`
	MiddleName string `json:"middleName" example:"Mejia"`
}

type PassengerTicketWholeInfo struct {
	Passenger
	Documents []DocumentTicketWholeInfo `json:"documents,omitempty"`
}
