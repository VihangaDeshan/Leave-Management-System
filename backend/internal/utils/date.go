package utils

import (
	"time"
)

// CalculateWorkingDays calculates the number of working days between two dates
// excluding weekends (Saturday and Sunday)
func CalculateWorkingDays(startDate, endDate time.Time) float64 {
	if endDate.Before(startDate) {
		return 0
	}

	// Normalize to start of day
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	workingDays := 0.0
	current := start

	for current.Before(end) || current.Equal(end) {
		// Check if current day is a weekday (Monday-Friday)
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			workingDays++
		}
		current = current.AddDate(0, 0, 1)
	}

	return workingDays
}

// IsWeekend checks if a given date is a weekend
func IsWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWorkingDay checks if a given date is a working day
func IsWorkingDay(date time.Time) bool {
	return !IsWeekend(date)
}

// GetStartOfDay returns the start of the day for a given time
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay returns the end of the day for a given time
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsDateRangeValid checks if end date is after or equal to start date
func IsDateRangeValid(startDate, endDate time.Time) bool {
	start := GetStartOfDay(startDate)
	end := GetStartOfDay(endDate)
	return !end.Before(start)
}

// HasDateOverlap checks if two date ranges overlap
func HasDateOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1) ||
		start1.Equal(start2) || end1.Equal(end2)
}

// FormatDateString formats a time.Time into a simple date string (YYYY-MM-DD)
func FormatDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseDateString parses a date string (YYYY-MM-DD) into time.Time
func ParseDateString(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// GetCurrentYear returns the current year
func GetCurrentYear() int {
	return time.Now().Year()
}
