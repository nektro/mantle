package db

import (
	"database/sql"

	"github.com/nektro/go-util/util"
	dbstorage "github.com/nektro/go.dbstorage"
)

type Role struct {
	ID                 int64  `json:"id"`
	UUID               string `json:"uuid" sqlite:"text"`
	Position           int    `json:"position" sqlite:"int"`
	Name               string `json:"name" sqlite:"text"`
	Color              string `json:"color" sqlite:"text"`
	PermManageChannels uint8  `json:"perm_manage_channels" sqlite:"tinyint(1)"`
	PermManageRoles    uint8  `json:"perm_manage_roles" sqlite:"tinyint(1)"`
}

//
//

func CreateRole(name string) string {
	id := db.QueryNextID(cTableRoles)
	uid := newUUID()
	util.Log("[role-create]", uid, name)
	db.Build().Ins(cTableRoles, id, uid, id, name, "", 1, 1).Exe()
	return uid
}

//
//

func (v Role) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(v.ID, v.UUID, v.Position, v.Name, v.Color, v.PermManageChannels, v.PermManageRoles)
	return &v
}

func (v Role) All() []*Role {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableRoles).Or("position", "asc"), Role{})
	res := []*Role{}
	for _, item := range arr {
		res = append(res, item.(*Role))
	}
	return res
}
