package internal

import (
	"github.com/linusback/aoc-go/internal/year2024"
	"github.com/linusback/aoc-go/pkg/errorsx"
)

func Solve(year, day string) (solution1, solution2 string, err error) {
	switch year {
	case "2024":
		return year2024.Solve(day)
	default:
		err = errorsx.NewSolverError(year, day, errorsx.ErrYearNotCreated)
		return
	}
}
