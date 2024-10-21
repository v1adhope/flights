package entities

type ReportRowByPassengerForPeriod struct {
	DateOfIssue     string `json:"dateOfIssue"`
	FlyAt           string `json:"flyAt"`
	TicketId        string `json:"ticketID"`
	FlyFrom         string `json:"flyFrom"`
	FlyTo           string `json:"flyTo"`
	ServiceProvided bool   `json:"serviceProvided"`
}

// type ReportFilterByPassengerForPeriod struct {
// 	PeriodFilter
// 	Kind bool
// }
