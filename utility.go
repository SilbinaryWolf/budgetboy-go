package main

import (
	"strconv"
	"time"
)

func TimeBeginningOfWeek(t time.Time, bSundayFirst bool) time.Time {

	weekday := int(t.Weekday())
	if !bSundayFirst {
		if weekday == 0 {
			weekday = 7
		}
		weekday = weekday - 1
	}

	d := time.Duration(-weekday) * 24 * time.Hour
	t = t.Add(d)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// TimeEndOfWeek return the end of the week of t
// bSundayFirst means that many country use the monday as the first day of week
func TimeEndOfWeek(t time.Time, bSundayFirst bool) time.Time {
	return TimeBeginningOfWeek(t, bSundayFirst).AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func DayOrdinal(x int) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return strconv.Itoa(x) + suffix
}
