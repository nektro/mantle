package db

import (
	"database/sql/driver"
	"errors"
	"strings"
)

// Array is a custom sql/driver type to handle list columns
type Array []string

// Scan - Implement the database/sql/driver Scanner interface
func (a *Array) Scan(value interface{}) error {
	if value == nil {
		*a = Array([]string{})
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			if len(v) > 0 {
				*a = Array(strings.Split(v, ","))
			} else {
				*a = Array([]string{})
			}
			return nil
		}
	}
	return errors.New("failed to scan Array")
}

// Value - Implement the database/sql/driver Valuer interface
func (a Array) Value() (driver.Value, error) {
	return strings.Join(a, ","), nil
}

func (a *Array) String() string {
	v, _ := a.Value()
	return v.(string)
}
