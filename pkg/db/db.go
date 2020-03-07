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
	db dbstorage.Database
)
var (
	Props = Properties{}
)

func Init() {
	db = etc.Database

	// table init
	db.CreateTableStruct(cTableSettings, Setting{})
	db.CreateTableStruct(cTableUsers, User{})
	db.CreateTableStruct(cTableChannels, Channel{})
	db.CreateTableStruct(cTableRoles, Role{})
	db.CreateTableStruct(cTableChannelPerms, ChannelPerm{})

	// load server properties
	Props.SetDefault("name", idata.Name)
	Props.SetDefault("owner", "")
	Props.SetDefault("public", "true")
	Props.SetDefault("description", "The new easy and effective communication platform for any successful team or community, providing you the messaging platform that puts you in charge of both the conversation and the data.")
	Props.SetDefault("cover_photo", "https://www.transparenttextures.com/patterns/gplay.png")
	Props.Init()
	Props.Set("version", idata.Version)
}

func Close() {
	// close db
	db.Close()
}
