package main

import (
	"net/http"
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

// Version takes in version string from build_all.sh
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
		item.AssertMessageTableExists()
	}

	//
	// add default channel, if none exist

	if len(_chans) == 0 {
		db.CreateChannel("general")
	}

	//
	// create server 'Owner' Role
	//		uneditable, and has all perms always

	pa := uint8(itypes.PermAllow)
	ws.RoleCache["o"] = &db.Role{
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
		db.Close()

		util.Log("Closing all remaining active WebSocket connections")
		ws.Close()

		util.Log("Done")
		os.Exit(0)
	})

	//
	// create http service

	r := etc.Router

	r.Path("/invite").Methods(http.MethodGet).HandlerFunc(handler.Invite)

	r.Path("/api/about").Methods(http.MethodGet).HandlerFunc(handler.ApiAbout)

	r.Path("/api/users/@me").Methods(http.MethodGet).HandlerFunc(handler.UsersMe)
	r.Path("/api/users/online").Methods(http.MethodGet).HandlerFunc(handler.UsersOnline)
	r.Path("/api/users/{uuid}").Methods(http.MethodGet).HandlerFunc(handler.UsersRead)

	r.Path("/api/channels/@me").Methods(http.MethodGet).HandlerFunc(handler.ChannelsMe)
	r.Path("/api/channels/create").Methods(http.MethodGet).HandlerFunc(handler.ChannelCreate)
	r.Path("/api/channels/{uuid}").Methods(http.MethodGet).HandlerFunc(handler.ChannelRead)
	r.Path("/api/channels/{uuid}/messages").Methods(http.MethodGet).HandlerFunc(handler.ChannelMessages)

	r.HandleFunc("/ws", handler.Websocket)

	//
	// start server

	etc.StartServer(idata.Config.Port)
}
