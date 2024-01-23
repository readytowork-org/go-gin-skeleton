package constants

import (
	"fmt"
)

type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
	Other  Gender = "O"
)

func (r Gender) IsValidVal(val string) error {
	switch val {
	case "M":
		return nil
	case "F":
		return nil
	case "O":
		return nil
	default:
		return fmt.Errorf("invalid user relation status: %q", val)
	}
}
