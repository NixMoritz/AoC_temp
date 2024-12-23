package util

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func GetYearDays(args []string) (year string, days []string, err error) {
	if len(args) > 2 {
		year = args[1]
		days = []string{args[2]}
		err = hasPassed(year, days[0])
		if err != nil {
			return "", nil, err
		}
	} else {
		year, days, err = getYearDays(time.Now())
		if err != nil {
			return "", nil, err
		}
	}
	return year, days, nil
}

func getYearDays(current time.Time) (year string, days []string, err error) {
	var loc *time.Location

	loc, err = time.LoadLocation("EST")
	if err != nil {
		return year, days, fmt.Errorf("while parsing location: %v", err)

	}

	current = time.Date(current.Year(), current.Month(), current.Day(), current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), loc)
	start := time.Date(current.Year(), time.November, 30, 0, 0, 0, 0, loc)
	y := current.Year()
	if current.Sub(start) < 0 {
		y -= 1
	}
	daysDiff := int64(current.Sub(start) / (24 * time.Hour))
	if daysDiff > 25 {
		daysDiff = 25
	}
	days = make([]string, 0, daysDiff)
	for i := daysDiff; i > 0; i-- {
		days = append(days, strconv.FormatInt(i, 10))
	}

	return strconv.Itoa(y), days, nil
}

func GetYearDay(args []string) (year string, days []string, err error) {
	var day string
	if len(args) > 2 {
		year = args[1]
		day = args[2]
		err = hasPassed(year, day)
		if err != nil {
			return "", nil, err
		}
		days = []string{day}
	} else if len(args) == 2 && args[1] == "all" {
		year, days, err = getYearDays(time.Now())
		if err != nil {
			return "", nil, err
		}
	} else {
		year, day, err = getYearDay(time.Now())
		if err != nil {
			return "", nil, err
		}
		days = []string{day}
	}
	return year, days, nil
}

func getYearDay(current time.Time) (year, day string, err error) {
	var loc *time.Location

	loc, err = time.LoadLocation("EST")
	if err != nil {
		return year, day, fmt.Errorf("while parsing location: %v", err)

	}

	current = time.Date(current.Year(), current.Month(), current.Day(), current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), loc)
	start := time.Date(current.Year(), time.November, 30, 0, 0, 0, 0, loc)
	y := current.Year()
	if current.Sub(start) < 0 {
		y -= 1
	}
	daysDiff := int64(current.Sub(start) / (24 * time.Hour))
	if daysDiff > 25 {
		daysDiff = 25
	}
	return strconv.Itoa(y), strconv.FormatInt(daysDiff, 10), nil
}

func hasPassed(year, day string) (err error) {
	var i int
	_, err = strconv.Atoi(year)
	if err != nil {
		return fmt.Errorf("while parsing year string: %v", err)
	}
	i, err = strconv.Atoi(day)
	if err != nil {
		return fmt.Errorf("while parsing day string: %v", err)
	}
	if i < 1 || i > 25 {
		return errors.New("day need to have a value between 1 and 25 inclusive")
	}
	return nil
}
