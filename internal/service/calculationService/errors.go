package calculationService

import "errors"

var (
	UnableToSaveCalculationError = errors.New("Unable to save calculation")
	InternalError                = errors.New("Internal error")
)
