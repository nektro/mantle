package main

import (
	"net/http"
	"os"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler"
	"github.com/nektro/mantle/pkg/iconst"
	"github.com/nektro/mantle/pkg/idata"
	"github.com/nektro/mantle/pkg/itypes"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	"github.com/spf13/pflag"

	_ "github.com/nektro/mantle/statik"
)

func main() {
	util.Log("Welcome to " + iconst.Name + " " + iconst.Version + ".")

	//
	pflag.IntVar(&idata.Config.Port, "port", 8000, "The port to bind the web server to.")
	etc.AppID = "mantle"
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
	// initialize server properties

	db.Props.SetDefault("name", iconst.Name)
	db.Props.SetDefault("owner", "")
	db.Props.SetDefault("public", "true")
	db.Props.Init()

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

	http.HandleFunc("/invite", handler.Invite)

	http.HandleFunc("/api/about", handler.ApiAbout)

	http.HandleFunc("/api/users/@me", handler.UsersMe)
	http.HandleFunc("/api/users/", handler.UsersRead)
	http.HandleFunc("/api/users/online", handler.UsersOnline)

	http.HandleFunc("/api/channels/@me", handler.ChannelsMe)
	http.HandleFunc("/api/channels/create", handler.ChannelCreate)

	http.HandleFunc("/ws", handler.Websocket)

	//
	// start server

	etc.StartServer(idata.Config.Port)
}
