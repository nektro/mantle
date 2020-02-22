package db

import (
	"github.com/nektro/mantle/pkg/iconst"
	"github.com/nektro/mantle/pkg/itypes"

	dbstorage "github.com/nektro/go.dbstorage"
	etc "github.com/nektro/go.etc"
)

var (
	DB dbstorage.Database
)

func Init() {
	DB = etc.Database

	// table init
	DB.CreateTableStruct(iconst.TableSettings, itypes.Setting{})
	DB.CreateTableStruct(iconst.TableUsers, itypes.User{})
	DB.CreateTableStruct(iconst.TableChannels, itypes.Channel{})
	DB.CreateTableStruct(iconst.TableRoles, itypes.Role{})
	DB.CreateTableStruct(iconst.TableChannelPerms, itypes.ChannelPerms{})
}
