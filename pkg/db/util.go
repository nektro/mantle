package db

import (
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid"
)

var epoch, _ = time.Parse("Jan 2 2006", "Jan 1 2020")

func newUUID() string {
	t := time.Unix(0, time.Now().UnixNano()-epoch.UnixNano())
	var entropy = ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func IsUID(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

func sUTCto3339(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.Replace(s, " ", "T", 1) + "Z"
}
