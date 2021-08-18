package asol

import (
	"fmt"
)

type NoRegisteredEventError struct {
}

func (error *NoRegisteredEventError) Error() string {
	return fmt.Sprintf("No event(s) registered.")
}

type ProcessNotFoundError struct {
	Process string
}

func (error *ProcessNotFoundError) Error() string {
	return fmt.Sprintf("%s could not be found", error.Process)
}

type TimeoutError struct {
	Process string
}

func (error *TimeoutError) Error() string {
	return fmt.Sprintf("%s could not be found and timed out", error.Process)
}
