package entities

import "time"

type Ticket struct {
	Id        string
	Provider  string
	FlyFrom   string
	FlyTo     string
	FlyAt     time.Time
	ArriveAt  time.Time
	CreatedAt time.Time
}
