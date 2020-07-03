package db

import (
	"database/sql"

	dbstorage "github.com/nektro/go.dbstorage"
)

type Setting struct {
	ID    int64  `json:"id"`
	Key   string `json:"key" dbsorm:"1"`
	Value string `json:"value" dbsorm:"1"`
}

//

func QuerySettingByKey(key string) (*Setting, bool) {
	ds, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableSettings).Wh("key", key), Setting{}).(*Setting)
	return ds, ok
}

//

// Scan implements dbstorage.Scannable
func (v Setting) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.Key, &v.Value)
	return &v
}

func (v Setting) All() []*Setting {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableSettings), Setting{})
	res := []*Setting{}
	for _, item := range arr {
		res = append(res, item.(*Setting))
	}
	return res
}
