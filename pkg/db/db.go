package db

import (
	"github.com/nektro/mantle/pkg/iconst"
	"github.com/nektro/mantle/pkg/idata"

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

	// load server properties
	Props.SetDefault("name", idata.Name)
	Props.SetDefault("owner", "")
	Props.SetDefault("public", "true")
	Props.Init()
}
