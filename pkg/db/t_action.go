package db

import (
	"database/sql/driver"
	"errors"
	"strconv"
)

// Action is a custom sql/driver type to handle Audit Log actions
type Action int16

// Actions Enum
const (
	ActionNull Action = iota
	ActionSettingUpdate
	ActionUserUpdate
)

// Value - Implement the database/sql Valuer interface
func (p Action) Value() (driver.Value, error) {
	return int64(p), nil
}

// Scan - Implement the database/sql Scanner interface
func (p *Action) Scan(value interface{}) error {
	if value == nil {
		*p = ActionNull
		return nil
	}
	if bv, err := driver.Int32.ConvertValue(value); err == nil {
		if v, ok := bv.(int64); ok {
			*p = Action(v)
			return nil
		}
	}
	return errors.New("failed to scan Action")
}

// String - Implement the fmt Stringer interface
func (p Action) String() string {
	return strconv.Itoa(int(p))
}
