package internal

import (
	"github.com/linusback/aoc-go-template/pkg/errorsx"
)

func Solve(year, day string) (solution1, solution2 string, err error) {
	switch year {
	default:
		err = errorsx.NewSolverError(year, day, errorsx.ErrYearNotCreated)
		return
	}
}
