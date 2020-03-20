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
	cTableInvites        = "invites"
)

var (
	db dbstorage.Database

	pa = PermAllow
	pi = PermIgnore
)
var (
	Props        = Properties{}
	BuiltInRoles = map[string]*Role{
		"o": &Role{
			0, "o", 0, "Owner", "", pa, pa, false, pa, pa, Time(epoch),
		},
	}
)

// Init sets up db tables and properties
func Init() {
	db = etc.Database

	// table init
	db.CreateTableStruct(cTableSettings, Setting{})
	db.CreateTableStruct(cTableUsers, User{})
	db.CreateTableStruct(cTableChannels, Channel{})
	db.CreateTableStruct(cTableRoles, Role{})
	db.CreateTableStruct(cTableInvites, Invite{})
	db.CreateTableStruct(cTableChannelPerms, ChannelPerm{})

	// load server properties
	Props.SetDefault("name", idata.Name)
	Props.SetDefault("owner", "")
	Props.SetDefault("public", "1")
	Props.SetDefault("description", "The new easy and effective communication platform for any successful team or community that's independently hosted and puts users, privacy, and effiecency first.")
	Props.SetDefault("cover_photo", "https://www.transparenttextures.com/patterns/gplay.png")
	Props.SetDefault("profile_photo", "https://avatars.discourse.org/v4/letter/m/ec9cab/90.png")
	Props.Init()
	Props.Set("version", idata.Version)

	// for loop create channel message tables
	_chans := (Channel{}.All())
	for _, item := range _chans {
		item.AssertMessageTableExists()
	}

	// add default channel, if none exist
	if len(_chans) == 0 {
		CreateChannel("general")
	}
}

// Close db
func Close() {
	db.Close()
}
