package db

import (
	"github.com/nektro/mantle/pkg/iconst"

	dbstorage "github.com/nektro/go.dbstorage"
	etc "github.com/nektro/go.etc"
)

var (
	DB    dbstorage.Database
	Props = Properties{}
)

func Init() {
	DB = etc.Database

	// table init
	DB.CreateTableStruct(iconst.TableSettings, Setting{})
	DB.CreateTableStruct(iconst.TableUsers, User{})
	DB.CreateTableStruct(iconst.TableChannels, Channel{})
	DB.CreateTableStruct(iconst.TableRoles, Role{})
	DB.CreateTableStruct(iconst.TableChannelPerms, ChannelPerms{})
}
