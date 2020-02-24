package db

import (
	"database/sql"

	dbstorage "github.com/nektro/go.dbstorage"
)

type Setting struct {
	ID    int64  `json:"id"`
	Key   string `json:"key" sqlite:"text"`
	Value string `json:"value" sqlite:"text"`
}

//

func QuerySettingByKey(key string) *Setting {
	rows := db.Build().Se("*").Fr(cTableSettings).Wh("key", key).Exe()
	defer rows.Close()
	if !rows.Next() {
		return nil
	}
	return Setting{}.Scan(rows).(*Setting)
}

//

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
