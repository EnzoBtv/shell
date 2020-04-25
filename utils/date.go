package utils

import "time"

// CommitDays defines the number of days of the last 6 months
const CommitDays = 183

// OutOfRange Control Const
const OutOfRange = 99999

// CommitWeeks defines the number of weeks of the last 6 months
const CommitWeeks = 26

// CalcOffset determines and returns the amount of days missing to fill
// the last row of the stats graph
func CalcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}

	return offset
}

// GetBeginningOfDay given a time.Time calculates the start time of that day
func GetBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

// CountDaysSinceDate counts how many days passed since the passed `date`
func CountDaysSinceDate(date time.Time) int {
	days := 0
	now := GetBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > CommitDays {
			return OutOfRange
		}
	}
	return days
}
