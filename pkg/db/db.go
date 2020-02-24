package db

import (
	"github.com/nektro/mantle/pkg/idata"

	dbstorage "github.com/nektro/go.dbstorage"
	etc "github.com/nektro/go.etc"
)

const (
	cTableSettings       = "server_settings"
	cTableUsers          = "users"
	cTableChannels       = "channels"
	cTableRoles          = "roles"
	cTableChannelPerms   = "channel_perms"
	cTableMessagesPrefix = "channel_messages_"
)

var (
	DB dbstorage.Database
)
var (
	Props = Properties{}
)

func Init() {
	DB = etc.Database

	// table init
	DB.CreateTableStruct(cTableSettings, Setting{})
	DB.CreateTableStruct(cTableUsers, User{})
	DB.CreateTableStruct(cTableChannels, Channel{})
	DB.CreateTableStruct(cTableRoles, Role{})
	DB.CreateTableStruct(cTableChannelPerms, ChannelPerms{})

	// load server properties
	Props.SetDefault("name", idata.Name)
	Props.SetDefault("owner", "")
	Props.SetDefault("public", "true")
	Props.Init()
	Props.Set("version", idata.Version)
}

func Close() {
	// close db
	DB.Close()
}
