package main

import (
	"github.com/linusback/aoc-go-template/pkg/aoc"
	"github.com/linusback/aoc-go-template/pkg/generate"
	"github.com/linusback/aoc-go-template/pkg/util"
	"log"
	"os"
	"slices"
	"time"
)

func main() {
	start := time.Now()
	log.Printf("Generating files\n\n")
	log.SetFlags(0)
	log.SetPrefix("  ")
	defer func() {
		log.SetPrefix("")
		log.Println("")
		log.SetFlags(log.LstdFlags)
		log.Printf("Time generating files: %v\n\n", time.Since(start))
	}()
	year, days, err := util.GetYearDays(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	slices.Reverse(days)
	err = generate.Generate(year, days)
	if err != nil {
		log.Fatal(err)
	}

	err = aoc.Download(year, days)
	if err != nil {
		log.Fatal(err)
	}
}
