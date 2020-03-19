package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler"
	"github.com/nektro/mantle/pkg/idata"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	"github.com/spf13/pflag"

	_ "github.com/nektro/mantle/statik"
)

// Version takes in version string from build_all.sh
var Version = "vMASTER"

func main() {
	idata.Version = etc.FixBareVersion(Version)
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

	r.Path("/").HandlerFunc(handler.InviteGet)
	r.Path("/invite").HandlerFunc(handler.InvitePost)
	r.Path("/verify").HandlerFunc(handler.Verify)

	r1 := r.PathPrefix("/api").Subrouter()
	r1.Path("/about").HandlerFunc(handler.ApiAbout)
	r1.Path("/update_property").HandlerFunc(handler.ApiPropertyUpdate)

	r2 := r1.PathPrefix("/users").Subrouter()
	r2.Path("/@me").HandlerFunc(handler.UsersMe)
	r2.Path("/online").HandlerFunc(handler.UsersOnline)
	r2.Path("/{uuid}").HandlerFunc(handler.UsersRead)
	r2.Path("/{uuid}/update").HandlerFunc(handler.UserUpdate)

	r3 := r1.PathPrefix("/channels").Subrouter()
	r3.Path("/@me").HandlerFunc(handler.ChannelsMe)
	r3.Path("/create").HandlerFunc(handler.ChannelCreate)
	r3.Path("/{uuid}").HandlerFunc(handler.ChannelRead)
	r3m := r3.Path("/{uuid}/messages").Subrouter()
	r3m.Methods(http.MethodGet).HandlerFunc(handler.ChannelMessagesRead)
	r3m.Methods(http.MethodDelete).HandlerFunc(handler.ChannelMessagesDelete)
	r3.Path("/{uuid}/update").HandlerFunc(handler.ChannelUpdate)

	r4 := r1.PathPrefix("/etc").Subrouter()
	r4b := r4.PathPrefix("/badges").Subrouter()
	r4b.Path("/members_online.svg").HandlerFunc(handler.EtcBadgeMembersOnline)
	r4b.Path("/members_total.svg").HandlerFunc(handler.EtcBadgeMembersTotal)
	r4.Path("/role_colors.css").HandlerFunc(handler.EtcRoleColorCSS)

	r5 := r1.PathPrefix("/roles").Subrouter()
	r5.Path("/@me").HandlerFunc(handler.RolesMe)
	r5.Path("/create").HandlerFunc(handler.RolesCreate)
	r5.Path("/{uuid}/update").HandlerFunc(handler.RoleUpdate)

	r6 := r1.PathPrefix("/invites").Subrouter()
	r6.Path("/@me").HandlerFunc(handler.InvitesMe)
	r6.Path("/create").HandlerFunc(handler.InvitesCreate)
	r.HandleFunc("/ws", handler.Websocket)

	//
	// start server

	etc.StartServer(idata.Config.Port)
}
