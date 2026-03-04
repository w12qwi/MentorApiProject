package postgres

import "errors"

var (
	UnableToConnectToPostgresError = errors.New("Unable to connect to postgres")
	NoSuchCalculationError         = errors.New("No such calculation")
)
