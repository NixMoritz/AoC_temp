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

func send(part aoc.Part, year, day, solution string) error {
	if len(solution) == 0 {
		log.Printf("empty answer part %s, skipping...\n", part)
		return nil
	}
	err := aoc.Send(part, year, day, solution)
	if err != nil {
		return fmt.Errorf("error while sending solution \"%s\" to probblem %s of day %s: %v", solution, part, day, err)
	}
	return nil
}
