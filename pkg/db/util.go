package db

import (
	"math/rand"
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

func now() Time {
	s := time.Now().UTC().String()[0:19]
	t, _ := time.Parse(timeFormat, s)
	return Time(t)
}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func uHighLow(a, b int) (int, int) {
	if a >= b {
		return a, b
	}
	return b, a
}
