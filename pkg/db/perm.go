package db

type Perm uint8

const (
	PermDeny Perm = iota
	PermIgnore
	PermAllow
)

// ToBool returns true if and only if this Perm is equal to PermAllow
func (p Perm) ToBool() bool {
	if p == PermAllow {
		return true
	}
	return false
}
