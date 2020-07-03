package db

import (
	"database/sql"

	dbstorage "github.com/nektro/go.dbstorage"
)

type ChannelPerm struct {
	ID        int64  `json:"id"`
	Channel   string `json:"channel" dbsorm:"1"`
	Type      int    `json:"p_type" dbsorm:"1"`
	Snowflake string `json:"snowflake" dbsorm:"1"`
}

// Scan implements dbstorage.Scannable
func (v ChannelPerm) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.Channel, &v.Type, &v.Snowflake)
	return &v
}

func (v ChannelPerm) All() []*ChannelPerm {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableChannelPerms), ChannelPerm{})
	res := []*ChannelPerm{}
	for _, item := range arr {
		res = append(res, item.(*ChannelPerm))
	}
	return res
}
