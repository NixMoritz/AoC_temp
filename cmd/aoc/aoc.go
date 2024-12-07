package main

import (
	"fmt"
	"github.com/linusback/aoc-go-template/internal"
	"github.com/linusback/aoc-go-template/pkg/aoc"
	"github.com/linusback/aoc-go-template/pkg/util"
	"log"
	"os"
	"time"
)

func main() {
	var (
		day  string
		err  error
		year string
	)

	log.Println("Start solver")
	year, day, err = util.GetYearDay(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var solution1, solution2 string
	start := time.Now()
	solution1, solution2, err = internal.Solve(year, day)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("result 1: ", solution1)
	log.Println("result 2: ", solution2)
	log.Printf("Time elapsed: %v\n\n", time.Since(start))
	log.Printf("Sending Answers")

	err = send(aoc.Part1, year, day, solution1)
	if err != nil {
		log.Println(err)
	}
	err = send(aoc.Part2, year, day, solution2)
	if err != nil {
		log.Println(err)
	}
}

// send wraps the return error from aoc.Send for more debuging information
func send(part aoc.Part, year, day, answer string) error {
	err := aoc.Send(part, year, day, answer)
	if err != nil {
		return fmt.Errorf("err for part %s of day %s: %v", part, day, err)
	}
	return nil
}
