package models

import "github.com/google/uuid"

// BINARY16 -> new datatype
type BINARY16 uuid.UUID

// StringToBinary16 -> parse string to id
func StringToBinary16(s string) (BINARY16, error) {
	id, err := uuid.Parse(s)
	return BINARY16(id), err
}

// String -> String Representation of Binary16
func (binary BINARY16) String() string {
	return uuid.UUID(binary).String()
}
