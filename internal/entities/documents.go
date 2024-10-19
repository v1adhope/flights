package entities

type Document struct {
	Id          string `json:"id" example:"uuid"`
	Type        string `json:"type" example:"Passport"`
	Number      string `json:"number" example:"3333777111"`
	PassengerId string `json:"passengerId" example:"uuid"`
}
