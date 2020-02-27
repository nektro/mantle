package db

import (
	"database/sql"

	dbstorage "github.com/nektro/go.dbstorage"
)

type ChannelPerm struct {
	ID        int64  `json:"id"`
	Channel   string `json:"channel" sqlite:"text"`
	Type      int    `json:"p_type" sqlite:"int"`
	Snowflake string `json:"snowflake" sqlite:"text"`
}

// Scan implements dbstorage.Scannable
func (v ChannelPerm) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.Channel, &v.Type, &v.Snowflake)
	return &v
}
