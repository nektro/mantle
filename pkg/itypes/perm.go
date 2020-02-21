package itypes

type Perm uint8

const (
	PermDeny Perm = iota
	PermIgnore
	PermAllow
)

func (p Perm) ToBool() bool {
	if p == PermAllow {
		return true
	}
	return false
}
