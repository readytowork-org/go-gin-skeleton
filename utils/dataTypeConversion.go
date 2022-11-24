package utils

import (
	"fmt"
	"strconv"
)

func Int64ToString(value int64) string {
	return fmt.Sprintf("%v", value)

}

func StringToInt64(value string) (*int64, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &i, err

}
