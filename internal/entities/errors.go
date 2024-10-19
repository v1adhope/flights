package entities

import "errors"

var (
	ErrorNothingToChange        = errors.New("Nothing to change")
	ErrorNothingToDelete        = errors.New("Nothing to delete")
	ErrorNothingFound           = errors.New("Nothing found")
	ErrorHasAlreadyExists       = errors.New("Has already exists")
	ErrorPassengerDoesNotExists = errors.New("Passenger doesn't exist")
)
