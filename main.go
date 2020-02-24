package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler"
	"github.com/nektro/mantle/pkg/idata"
	"github.com/nektro/mantle/pkg/itypes"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	"github.com/spf13/pflag"

	_ "github.com/nektro/mantle/statik"
)

var Version = "vMASTER"

func main() {
	idata.Version = etc.FixBareVersion(Version) + "-" + runtime.Version()
	util.Log("Welcome to " + idata.Name + " " + idata.Version + ".")

	//
	pflag.IntVar(&idata.Config.Port, "port", 8000, "The port to bind the web server to.")
	etc.AppID = strings.ToLower(idata.Name)
	etc.PreInit()

	//
	etc.Init("mantle", &idata.Config, "./invite", handler.SaveOAuth2InfoCb)

	//
	// database initialization

	db.Init()

	// for loop create channel message tables
	_chans := (db.Channel{}.All())
	for _, item := range _chans {
		db.AssertChannelMessagesTableExists(item.UUID)
	}

	//
	// add default channel, if none exist

	if len(_chans) == 0 {
		db.CreateChannel("chat")
	}

	//
	// create server 'Owner' Role
	//		uneditable, and has all perms always

	pa := uint8(itypes.PermAllow)
	ws.RoleCache["o"] = db.Role{
		0, "o", 0, "Owner", "", pa, pa,
	}

	//
	// load roles into local cache

	for _, item := range (db.Role{}.All()) {
		ws.RoleCache[item.UUID] = item
	}

	//
	// setup graceful stop

	util.RunOnClose(func() {
		util.Log("Gracefully shutting down...")

		util.Log("Saving database to disk")
		db.DB.Close()

		util.Log("Closing all remaining active WebSocket connections")
		for _, item := range ws.UserCache {
			item.Conn.Close()
		}

		util.Log("Done")
		os.Exit(0)
	})

	//
	// create http service

	r := etc.Router

	r.HandleFunc("/invite", handler.Invite)

	r.HandleFunc("/api/about", handler.ApiAbout)

	r.HandleFunc("/api/users/@me", handler.UsersMe)
	r.HandleFunc("/api/users/online", handler.UsersOnline)
	r.HandleFunc("/api/users/{uuid:[0-9a-f]{32}}", handler.UsersRead)

	r.HandleFunc("/api/channels/@me", handler.ChannelsMe)
	r.HandleFunc("/api/channels/create", handler.ChannelCreate)

	r.HandleFunc("/ws", handler.Websocket)

	//
	// start server

	etc.StartServer(idata.Config.Port)
}
