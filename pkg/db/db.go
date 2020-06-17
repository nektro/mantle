package db

import (
	"strconv"

	"github.com/nektro/mantle/pkg/idata"

	"github.com/nektro/go-util/util"
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
	cTableAudits         = "audits"
)

var (
	// ResourceTables is the list of db table names that represent the various resources in Mantle
	ResourceTables = []string{cTableUsers, cTableChannels, cTableRoles, cTableInvites, cTableAudits}
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
			0, "o", 0, "Owner", "", pa, pa, false, pa, pa, Time(Epoch), pa,
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
	db.DropTable(cTableChannelPerms)
	db.CreateTableStruct(cTableAudits, Audit{})

	// load server properties
	Props.SetDefault("name", idata.Name)
	Props.SetDefault("owner", "")
	Props.SetDefault("public", "1")
	Props.SetDefault("description", "The new easy and effective communication platform for any successful team or community that's independently hosted and puts users, privacy, and effiecency first.")
	Props.SetDefault("cover_photo", "data:,")
	Props.SetDefault("profile_photo", "https://avatars.discourse.org/v4/letter/m/ec9cab/90.png")
	Props.SetDefault("prometheus_key", util.RandomString(64))

	Props.SetDefaultInt64("count_users_members", queryCount(db.Build().Se("*").Fr(cTableUsers).Wh("is_member", "1").Exe()))
	Props.SetDefaultInt64("count_users_banned", queryCount(db.Build().Se("*").Fr(cTableUsers).Wh("is_banned", "1").Exe()))
	Props.SetDefaultInt64("count_users_members_max", 0)

	for _, item := range ResourceTables {
		Props.SetDefaultInt64("count_"+item, db.QueryRowCount(item))
	}
	for i := 1; i < ActionLen(); i++ {
		is := strconv.Itoa(i)
		Props.SetDefaultInt64("count_"+cTableAudits+"_action_"+is, queryCount(db.Build().Se("*").Fr(cTableAudits).Wh("action", is).Exe()))
	}

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
