package entities

type Document struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Number      string `json:"number"`
	PassangerId string `json:"passangerId"`
}
