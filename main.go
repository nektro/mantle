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
	"github.com/nektro/go-util/vflag"
	etc "github.com/nektro/go.etc"
	"github.com/nektro/go.etc/htp"
	"github.com/nektro/go.etc/translations"

	_ "github.com/nektro/mantle/statik"
)

// Version takes in version string from build_all.sh
var Version = "vMASTER"

func main() {
	rand.Seed(time.Now().UnixNano())

	idata.Version = etc.FixBareVersion(Version)
	util.Log("Welcome to " + idata.Name + " " + idata.Version + ".")

	//
	vflag.IntVar(&idata.Config.Port, "port", 8000, "The port to bind the web server to.")
	etc.AppID = strings.ToLower(idata.Name)
	etc.PreInit()

	//
	etc.Init("mantle", &idata.Config, "./verify", handler.SaveOAuth2InfoCb)

	//
	// database initialization

	db.Init()

	translations.Fetch()
	translations.Init()

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

	fRegister("/", sPaths{
		GET: handler.InviteGet,
		Sub: map[string]sPaths{
			"invite":  sPaths{POS: handler.InvitePost},
			"verify":  sPaths{GET: handler.Verify},
			"ws":      sPaths{GET: handler.Websocket},
			"chat": sPaths{
				Sub: map[string]sPaths{
					"": sPaths{GET: handler.Chat},
				},
			},
			"api": sPaths{
				Sub: map[string]sPaths{
					"about":           sPaths{GET: handler.ApiAbout},
					"update_property": sPaths{PUT: handler.ApiPropertyUpdate},
					"etc": sPaths{
						Sub: map[string]sPaths{
							"role_colors.css": sPaths{GET: handler.EtcRoleColorCSS},
							"badges": sPaths{
								Sub: map[string]sPaths{
									"members_online.svg": sPaths{GET: handler.EtcBadgeMembersOnline},
									"members_total.svg":  sPaths{GET: handler.EtcBadgeMembersTotal},
								},
							},
						},
					},
					"users": sPaths{
						Sub: map[string]sPaths{
							"@me":    sPaths{GET: handler.UsersMe},
							"online": sPaths{GET: handler.UsersOnline},
							"{uuid}": sPaths{
								GET: handler.UsersRead,
								PUT: handler.UserUpdate,
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
					"roles": sPaths{
						Sub: map[string]sPaths{
							"@me":    sPaths{GET: handler.RolesMe},
							"create": sPaths{POS: handler.RolesCreate},
							"{uuid}": sPaths{
								PUT: handler.RoleUpdate,
								DEL: handler.RoleDelete,
							},
						},
					},
					"invites": sPaths{
						Sub: map[string]sPaths{
							"@me":    sPaths{GET: handler.InvitesMe},
							"create": sPaths{POS: handler.InvitesCreate},
							"{uuid}": sPaths{
								PUT: handler.InviteUpdate,
								DEL: handler.InviteDelete,
							},
						},
					},
					"admin": sPaths{
						Sub: map[string]sPaths{
							"audits.csv": sPaths{GET: handler.AuditsCsv},
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

func iregister(m, p string, h http.HandlerFunc) {
	if h == nil {
		return
	}
	htp.Register(m, p, h)
}

func fRegister(s string, p sPaths) {
	if strings.HasPrefix(s, "//") {
		s = s[1:]
	}
	iregister(http.MethodGet, s, p.GET)
	iregister(http.MethodPost, s, p.POS)
	iregister(http.MethodPut, s, p.PUT)
	iregister(http.MethodDelete, s, p.DEL)
	//
	if p.Sub != nil {
		for k, v := range p.Sub {
			fRegister(s+"/"+k, v)
		}
	}
}
