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
	etc.Init("mantle", &idata.Config, "./verify", handler.SaveOAuth2InfoCb)

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

	r.Path("/verify").HandlerFunc(handler.Verify)

	r1 := r.PathPrefix("/api").Subrouter()
	r1.Path("/about").HandlerFunc(handler.ApiAbout)

	r2 := r1.PathPrefix("/users").Subrouter()
	r2.Path("/@me").HandlerFunc(handler.UsersMe)
	r2.Path("/online").HandlerFunc(handler.UsersOnline)
	r2.Path("/{uuid}").HandlerFunc(handler.UsersRead)

	r3 := r1.PathPrefix("/channels").Subrouter()
	r3.Path("/@me").HandlerFunc(handler.ChannelsMe)
	r3.Path("/create").HandlerFunc(handler.ChannelCreate)
	r3.Path("/{uuid}").HandlerFunc(handler.ChannelRead)
	r3m := r3.Path("/{uuid}/messages").Subrouter()
	r3m.Methods(http.MethodGet).HandlerFunc(handler.ChannelMessagesRead)
	r3m.Methods(http.MethodDelete).HandlerFunc(handler.ChannelMessagesDelete)

	r4 := r1.PathPrefix("/etc").Subrouter()
	r4b := r4.PathPrefix("/badges").Subrouter()
	r4b.Path("/members_online.svg").HandlerFunc(handler.EtcBadgeMembersOnline)
	r4b.Path("/members_total.svg").HandlerFunc(handler.EtcBadgeMembersTotal)

	r.HandleFunc("/ws", handler.Websocket)

	//
	// start server

	etc.StartServer(idata.Config.Port)
}
