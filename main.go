package main

import (
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

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
	rand.Seed(time.Now().UnixNano())

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
	r2u := r2.PathPrefix("/{uuid}").Subrouter()
	r2u.Path("").HandlerFunc(handler.UsersRead)
	r2u.Methods(http.MethodPut).HandlerFunc(handler.UserUpdate)

	r3 := r1.PathPrefix("/channels").Subrouter()
	r3.Path("/@me").HandlerFunc(handler.ChannelsMe)
	r3.Path("/create").HandlerFunc(handler.ChannelCreate)
	r3u := r3.PathPrefix("/{uuid}").Subrouter()
	r3v := r3u.Path("")
	r3v.Methods(http.MethodGet).HandlerFunc(handler.ChannelRead)
	r3v.Methods(http.MethodPut).HandlerFunc(handler.ChannelUpdate)
	r3v.Methods(http.MethodDelete).HandlerFunc(handler.ChannelDelete)
	r3m := r3u.PathPrefix("/messages").Subrouter()
	r3m.Methods(http.MethodGet).HandlerFunc(handler.ChannelMessagesRead)
	r3m.Methods(http.MethodDelete).HandlerFunc(handler.ChannelMessagesDelete)

	r4 := r1.PathPrefix("/etc").Subrouter()
	r4b := r4.PathPrefix("/badges").Subrouter()
	r4b.Path("/members_online.svg").HandlerFunc(handler.EtcBadgeMembersOnline)
	r4b.Path("/members_total.svg").HandlerFunc(handler.EtcBadgeMembersTotal)
	r4.Path("/role_colors.css").HandlerFunc(handler.EtcRoleColorCSS)

	r5 := r1.PathPrefix("/roles").Subrouter()
	r5.Path("/@me").HandlerFunc(handler.RolesMe)
	r5.Path("/create").HandlerFunc(handler.RolesCreate)
	r5u := r5.PathPrefix("/{uuid}").Subrouter()
	r5u.Methods(http.MethodPut).HandlerFunc(handler.RoleUpdate)
	r5u.Methods(http.MethodDelete).HandlerFunc(handler.RoleDelete)

	r6 := r1.PathPrefix("/invites").Subrouter()
	r6.Path("/@me").HandlerFunc(handler.InvitesMe)
	r6.Path("/create").HandlerFunc(handler.InvitesCreate)
	r6u := r6.PathPrefix("/{uuid}").Subrouter()
	r6u.Methods(http.MethodPut).HandlerFunc(handler.InviteUpdate)
	r6u.Methods(http.MethodDelete).HandlerFunc(handler.InviteDelete)

	r.HandleFunc("/ws", handler.Websocket)
	fRegister("/", sPaths{
	})

	//
	// start server

	etc.StartServer(idata.Config.Port)
}

type sPaths struct {
	GET http.HandlerFunc
	POS http.HandlerFunc
	PUT http.HandlerFunc
	DEL http.HandlerFunc
	Sub map[string]sPaths
}

func fRegister(s string, p sPaths) {
	if strings.HasPrefix(s, "//") {
		s = s[1:]
	}
	if p.GET != nil {
		etc.Router.HandleFunc(s, p.GET)
	}
	if p.POS != nil {
		etc.Router.HandleFunc(s, p.POS)
	}
	if p.PUT != nil {
		etc.Router.HandleFunc(s, p.PUT)
	}
	if p.DEL != nil {
		etc.Router.HandleFunc(s, p.DEL)
	}
	if p.Sub != nil {
		for k, v := range p.Sub {
			fRegister(s+"/"+k, v)
		}
	}
}
