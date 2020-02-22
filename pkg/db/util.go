package db

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func newUUID() string {
	return strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
}
