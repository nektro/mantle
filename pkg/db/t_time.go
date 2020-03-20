package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// timeFormat is the underlying format of Time.String()[:19]
const timeFormat = "2006-01-02 15:04:05"

// timeZero is empty time
const timeZero = "0000-01-01 00:00:00"

// Time is a custom sql/driver type to handle UTC times
type Time time.Time

// NewTime accepts a string in timeFormat and returns a Time. timeZero on error.
func NewTime(s string) Time {
	r, err := time.Parse(timeFormat, s)
	if err != nil {
		NewTime(timeZero)
	}
	return Time(r)
}

// T returns the underlying time.Time object
func (t Time) T() time.Time {
	return time.Time(t)
}

// Scan - Implement the database/sql/driver Scanner interface
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		// *t = Time([]string{})
		return nil
	}
	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			if len(v) == 0 {
				v = timeZero
			}
			*t = NewTime(v)
			return nil
		}
	}
	return errors.New("failed to scan Time")
}

// Value - Implement the database/sql/driver Valuer interface
func (t Time) Value() (driver.Value, error) {
	s := t.T().String()
	if len(s) == 0 {
		return NewTime(timeZero), nil
	}
	return s[:19], nil
}

// String - Implement the fmt Stringer interface
func (t Time) String() string {
	v, _ := t.Value()
	s := v.(string)
	if s == timeZero {
		return ""
	}
	return strings.Replace(s, " ", "T", 1) + "Z"
}

// MarshalJSON - Implement the json Marshal interface
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
