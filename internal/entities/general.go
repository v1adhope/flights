package entities

type Id struct {
	Value string `json:"id" example:"uuid"`
}

type PeriodFilter struct {
	From string
	To   string
}
