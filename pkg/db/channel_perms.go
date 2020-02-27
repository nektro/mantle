package db

import (
	"database/sql"

	dbstorage "github.com/nektro/go.dbstorage"
)

type ChannelPerms struct {
	ID        int64  `json:"id"`
	Channel   string `json:"channel" sqlite:"text"`
	Type      int    `json:"p_type" sqlite:"int"`
	Snowflake string `json:"snowflake" sqlite:"text"`
}

func (v ChannelPerms) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.Channel, &v.Type, &v.Snowflake)
	return &v
}
