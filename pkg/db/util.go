package db

import (
	"database/sql"
	"time"

	. "github.com/nektro/go.etc/dbt"
)

func now() Time {
	s := time.Now().UTC().String()[0:19]
	t, _ := time.Parse(TimeFormat, s)
	return Time(t)
}

func uHighLow(a, b int) (int, int) {
	if a >= b {
		return a, b
	}
	return b, a
}

func queryCount(rows *sql.Rows) int64 {
	s := int64(0)
	defer rows.Close()
	for rows.Next() {
		s++
	}
	return s
}

type iIDers interface {
	t() string
	i() string
}

func doUp(v iIDers, col string, value string) {
	db.Build().Up(v.t(), col, value).Wh("uuid", v.i()).Exe()
}

func doDel(v iIDers) {
	db.Build().Del(v.t()).Wh("uuid", v.i()).Exe()
}
