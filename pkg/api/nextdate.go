package api

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("the repetition rule is empty")
	}
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("can't convert start date to correct date: %w", err)
	}
	repeatSlice := strings.Split(repeat, " ")
	switch repeatSlice[0] {
	case "y":
		if len(repeatSlice) > 1 {
			return "", fmt.Errorf("interval specified: %s", repeat)
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil
	case "d":
		if len(repeatSlice) == 1 {
			return "", fmt.Errorf("interval in days not specified: %s", repeat)
		}
		if len(repeatSlice) > 2 {
			return "", fmt.Errorf("too many values for d rule: %s", repeat)
		}
		days, err := strconv.Atoi(repeatSlice[1])
		if err != nil {
			return "", fmt.Errorf("conversion error to number: %w", err)
		}
		if days < 1 || days > 400 {
			return "", fmt.Errorf("invalid interval: %s", repeat)
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil
	case "w":
		if len(repeatSlice) == 1 {
			return "", fmt.Errorf("interval in w not specified: %s", repeat)
		}
		if len(repeatSlice) > 2 {
			return "", fmt.Errorf("too many values for w rule: %s", repeat)
		}
		weekdaysStr := strings.Split(repeatSlice[1], ",")
		if len(weekdaysStr) > 7 {
			return "", fmt.Errorf("too many days for w rule: %s", repeat)
		}
		weekdays := make([]int, 0, 7)
		for _, v := range weekdaysStr {
			dw, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("conversion error to number: %w", err)
			}
			if dw < 1 || dw > 7 {
				return "", fmt.Errorf("not a week number: %d", dw)
			} else {
				weekdays = append(weekdays, dw)
			}
		}
		if !afterNow(date, now) {
			date = now
		}
		weekdayDec := dayOfWeekNumber(date)
		sort.Ints(weekdays)
		for _, v := range weekdays {
			if v > weekdayDec {
				date = date.AddDate(0, 0, v-weekdayDec)
				return date.Format(DateFormat), nil
			}
		}
		date = date.AddDate(0, 0, 7+weekdays[0]-weekdayDec)
		return date.Format(DateFormat), nil
	case "m":
		var days, months []int
		if len(repeatSlice) < 2 || len(repeatSlice) > 3 {
			return "", fmt.Errorf("invalid m rule format")
		}
		days, err = parseNumbers(repeatSlice[1])
		if err != nil {
			return "", fmt.Errorf("invalid month days format")
		}
		if !isValidMonthDay(days) {
			return "", fmt.Errorf("month days must be between 1 and 31 or -1, -2")
		}
		if len(repeatSlice) == 3 {
			months, err = parseNumbers(repeatSlice[2])
			if err != nil {
				return "", fmt.Errorf("invalid months format")
			}
			if !isValidMonth(months) {
				return "", fmt.Errorf("months must be between 1 and 12")
			}
		} else {
			months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}
		for {
			date = date.AddDate(0, 0, 1)
			currentMonth := int(date.Month())
			currentDay := date.Day()
			for _, month := range months {
				if currentMonth != month {
					continue
				}
				for _, md := range days {
					if md > 0 && md == currentDay {
						if afterNow(date, now) {
							return date.Format(DateFormat), nil
						}
					} else if md < 0 {
						lastDay := getLastDayOfMonth(date)
						if currentDay == lastDay+md+1 {
							if afterNow(date, now) {
								return date.Format(DateFormat), nil
							}
						}
					}
				}
			}
			if date.After(now.AddDate(100, 0, 0)) {
				return "", fmt.Errorf("no valid month day found within reasonable time")
			}
		}
	default:
		return "", fmt.Errorf("invalid character: %s", repeat)
	}
}

func afterNow(date, now time.Time) bool {
	if date.After(now) {
		return true
	}
	return false
}

func dayOfWeekNumber(date time.Time) int {
	weekday := date.Weekday()
	if weekday == 0 {
		return 7
	}
	return int(weekday)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func parseNumbers(s string) ([]int, error) {
	if s == "" {
		return nil, fmt.Errorf("empty number list")
	}
	numStrs := strings.Split(s, ",")
	nums := make([]int, 0, len(numStrs))
	for _, numStr := range numStrs {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			return nil, fmt.Errorf("invalid number format")
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func isValidMonthDay(days []int) bool {
	for _, day := range days {
		if day < -2 || day == 0 || day > 31 {
			return false
		}
	}
	return true
}

func isValidMonth(months []int) bool {
	for _, month := range months {
		if month < 1 || month > 12 {
			return false
		}
	}
	return true
}

func getLastDayOfMonth(t time.Time) int {
	return t.AddDate(0, 1, -t.Day()).Day()
}
