package dates

import (
	"boilerplate-api/dtos"
	"time"
)

func LastTwoDates() dtos.Dates {
	currentDate := time.Now()
	currentMonthFirstDay := time.Date(
		currentDate.Year(),
		currentDate.Month(),
		1,
		0,
		0,
		0,
		0,
		currentDate.Location(),
	)
	return dtos.Dates{
		LastMonth: dtos.Date{
			FirstDay: currentMonthFirstDay.AddDate(0, -1, 0),
			LastDay:  currentMonthFirstDay.AddDate(0, 0, -1),
		},
		CurrentMonth: dtos.Date{
			FirstDay: currentMonthFirstDay,
			LastDay:  currentMonthFirstDay.AddDate(0, 1, -1),
		},
	}
}
