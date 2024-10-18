package entities

import "errors"

var (
	ErrorNothingToChange = errors.New("Nothing to change")
	ErrorNothingToDelete = errors.New("Nothing to delete")
	ErrorNothingFound    = errors.New("Nothing found")
)
