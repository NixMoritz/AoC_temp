package errorsx

import (
	"fmt"
)

const (
	ErrYearNotCreated SimpleError = "year not yet created"
	ErrDayNotCreated  SimpleError = "day not yet created"
)

type SimpleError string

func (se SimpleError) Error() string {
	return string(se)
}

type SolveError struct {
	year, day string
	err       error
}

func (se *SolveError) Error() string {
	return fmt.Sprintf("from solver year %s, day %s: %v", se.year, se.day, se.err)
}

func (se *SolveError) Unwrap() error {
	return se.err
}

func NewSolverError(year, day string, err error) error {
	if err == nil {
		return nil
	}
	return &SolveError{
		year: year,
		day:  day,
		err:  err,
	}
}
