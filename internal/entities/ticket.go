package entities

type Ticket struct {
	Id        string `json:"id"`
	Provider  string `json:"provider"`
	FlyFrom   string `json:"flyFrom"`
	FlyTo     string `json:"flyTo"`
	FlyAt     string `json:"flyAt"`
	ArriveAt  string `json:"arriveAt"`
	CreatedAt string `json:"createdAt"`
}
