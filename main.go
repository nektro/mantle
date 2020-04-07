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
	fRegister("/", sPaths{
		GET: handler.InviteGet,
		Sub: map[string]sPaths{
			"invite":  sPaths{POS: handler.InvitePost},
			"verify":  sPaths{GET: handler.Verify},
			"ws":      sPaths{GET: handler.Websocket},
			"api": sPaths{
				Sub: map[string]sPaths{
					"about":           sPaths{GET: handler.ApiAbout},
					"update_property": sPaths{PUT: handler.ApiPropertyUpdate},
					"users": sPaths{
						Sub: map[string]sPaths{
							"@me":    sPaths{GET: handler.UsersMe},
							"online": sPaths{GET: handler.UsersOnline},
							"{uuid}": sPaths{
								GET: handler.UsersRead,
							},
						},
					},
					"channels": sPaths{
						Sub: map[string]sPaths{
							"@me":    sPaths{GET: handler.ChannelsMe},
							"create": sPaths{POS: handler.ChannelCreate},
							"{uuid}": sPaths{
								GET: handler.ChannelRead,
								PUT: handler.ChannelUpdate,
								DEL: handler.ChannelDelete,
								Sub: map[string]sPaths{
									"messages": sPaths{
										GET: handler.ChannelMessagesRead,
										DEL: handler.ChannelMessagesDelete,
									},
								},
							},
						},
					},
			},
		},
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
