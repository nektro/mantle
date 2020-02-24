package db

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func newUUID() string {
	return strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
}
