package calculationService

import "errors"

var (
	UnableToSaveCalculationError = errors.New("Unable to save calculation")
	NoSuchCalculationError       = errors.New("No such calculation")
	InternalError                = errors.New("Internal error")
)
