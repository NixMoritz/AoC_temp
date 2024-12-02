package year2024

import (
	"github.com/linusback/aoc-go/internal/year2024/day1"
	"github.com/linusback/aoc-go/pkg/errorsx"
)

const year = "2024"

func Solve(day string) (solution1, solution2 string, err error) {
	switch day {
	case "1":
		return day1.Solve()
	default:
		err = errorsx.NewSolverError(year, day, errorsx.ErrDayNotCreated)
		return
	}
}
