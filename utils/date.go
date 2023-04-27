package utils

import (
	"fmt"
	"time"
)

// ConvertStringToDate Converts the string to date
func ConvertStringToDate(date string) time.Time {
	_time, err := time.Parse("2006-01-02", date)
	fmt.Println(err)
	return _time
}

// ConvertRFCStringToDate Convert RFC date string to date
func ConvertRFCStringToDate(date string) time.Time {
	_time, err := time.Parse(time.RFC3339, date)
	fmt.Println(err)
	return _time
}
