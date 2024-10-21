package v1

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var names validator.Func = func(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	isMatched, err := regexp.MatchString("^[A-z ,.'-]+$", value)
	if err != nil {
		return false
	}

	return isMatched
}

func ticketCreateReqStructLevelValidation(sl validator.StructLevel) {
	ticket := sl.Current().Interface().(ticketCreateReq)

	flyAt, err := time.Parse(time.RFC3339, ticket.FlyAt)
	if err != nil {
		sl.ReportError(ticket.FlyAt, "flyAt", "FlyAt", "rfc3339Time", "")
	}

	arriveAt, err := time.Parse(time.RFC3339, ticket.ArriveAt)
	if err != nil {
		sl.ReportError(ticket.ArriveAt, "arriveAt", "ArriveAt", "rfc3339Time", "")
	}

	difference := flyAt.Sub(time.Now())
	if difference < 0 {
		sl.ReportError(ticket.FlyAt, "flyAt", "FlyAt", "flyght_before_now", "")
	}

	difference = arriveAt.Sub(flyAt)
	if difference < 0 {
		sl.ReportError(ticket.ArriveAt, "arriveAt", "ArriveAt", "arrive_before_fly", "")
	}
}

func reportByPassengerIdForPeriodQueryStructLevelValidation(sl validator.StructLevel) {
	query := sl.Current().Interface().(reportByPassengerIdForPeriodQuery)

	from, err := time.Parse(time.RFC3339, query.From)
	if err != nil {
		sl.ReportError(query.From, "from", "From", "rfc3339Time", "")
	}

	to, err := time.Parse(time.RFC3339, query.To)
	if err != nil {
		sl.ReportError(query.To, "to", "To", "rfc3339Time", "")
	}

	difference := to.Sub(from)
	if difference < 0 {
		sl.ReportError(query.To, "to", "To", "from_after_to", "")
	}
}
