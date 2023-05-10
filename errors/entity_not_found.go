package errors

import "fmt"

type EntityNotFound struct {
	Entity string
	ID     string
}

func (e EntityNotFound) Error() string {
	return fmt.Sprintf("No %v found for ID: %v", e.Entity, e.ID)
}
