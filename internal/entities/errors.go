package entities

import "errors"

var (
	ErrorNothingToChange            = errors.New("Nothing to change")
	ErrorNothingToDelete            = errors.New("Nothing to delete")
	ErrorNothingFound               = errors.New("Nothing found")
	ErrorHasAlreadyExists           = errors.New("Has already exists")
	ErrorPassengerDoesNotExists     = errors.New("Passenger doesn't exist")
	ErrorTicketDoesNotExists        = errors.New("Ticket doesn't exist")
	ErrorsThereArePassengersOnRaice = errors.New("There are passengers on raice")
)
