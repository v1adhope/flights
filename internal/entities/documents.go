package entities

type Document struct {
	Id          string `json:"id" example:"uuid"`
	Type        string `json:"type" example:"Passport"`
	Number      string `json:"number" example:"3333777111"`
	PassengerId string `json:"passengerId" example:"uuid"`
}

type DocumentTicketWholeInfo struct {
	Id     string `json:"id,omitempty" example:"uuid"`
	Type   string `json:"type,omitempty" example:"Passport"`
	Number string `json:"number,omitempty" example:"3333777111"`
}
