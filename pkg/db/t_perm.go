package db

import (
	"database/sql/driver"
	"errors"
)

// Perm is a custom sql/driver type to handle 3-state permissions
type Perm uint8

const (
	PermDeny Perm = iota
	PermIgnore
	PermAllow
)

// ToBool returns true if and only if this Perm is equal to PermAllow
func (p Perm) ToBool() bool {
	return p == PermAllow
}

// Value - Implement the database/sql Valuer interface
func (p Perm) Value() (driver.Value, error) {
	return int64(p), nil
}

// Scan - Implement the database/sql Scanner interface
func (p *Perm) Scan(value interface{}) error {
	if value == nil {
		*p = PermIgnore
		return nil
	}
	if bv, err := driver.Int32.ConvertValue(value); err == nil {
		if v, ok := bv.(int64); ok {
			*p = Perm(v)
			return nil
		}
	}
	return errors.New("failed to scan Perm")
}
