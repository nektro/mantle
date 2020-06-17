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
	"github.com/nektro/mantle/pkg/metrics"
	"github.com/nektro/mantle/pkg/store"
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

	etc.AppID = strings.ToLower(idata.Name)
	idata.Version = etc.FixBareVersion(Version)
	idata.Version = strings.ReplaceAll(idata.Version, "-", ".")
	util.Log("Starting " + idata.Name + " " + idata.Version + ".")

	//
	vflag.StringVar(&idata.Config.RedisURL, "redis-url", "", "")
	vflag.IntVar(&idata.Config.MaxMemberCount, "max-member-count", 0, "")

	etc.PreInit()
	etc.Init(&idata.Config, "./verify", handler.SaveOAuth2InfoCb)

	//
	// database initialization

	store.Init()

	db.Init()

	translations.Fetch()
	translations.Init()

	metrics.Init()

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

	handler.Init()

	htp.Register("/", http.MethodGet, handler.InviteGet)
	htp.Register("/invite", http.MethodGet, handler.InvitePost)
	htp.Register("/verify", http.MethodGet, handler.Verify)
	htp.Register("/ws", http.MethodGet, handler.Websocket)
	htp.Register("/metrics", http.MethodGet, metrics.Handler())

	htp.Register("/chat/", http.MethodGet, handler.Chat)

	htp.Register("/api/about", http.MethodGet, handler.ApiAbout)
	htp.Register("/api/update_property", http.MethodPut, handler.ApiPropertyUpdate)

	htp.Register("/api/etc/role_colors.css", http.MethodGet, handler.EtcRoleColorCSS)

	htp.Register("/api/etc/badges/version.svg", http.MethodGet, handler.EtcBadgeVersion)
	htp.Register("/api/etc/badges/members_online.svg", http.MethodGet, handler.EtcBadgeMembersOnline)
	htp.Register("/api/etc/badges/members_total.svg", http.MethodGet, handler.EtcBadgeMembersTotal)

	htp.Register("/api/users/@me", http.MethodGet, handler.UsersMe)
	htp.Register("/api/users/online", http.MethodGet, handler.UsersOnline)
	htp.Register("/api/users/{uuid}", http.MethodGet, handler.UsersRead)
	htp.Register("/api/users/{uuid}", http.MethodPut, handler.UserUpdate)

	htp.Register("/api/channels/@me", http.MethodGet, handler.ChannelsMe)
	htp.Register("/api/channels/create", http.MethodPost, handler.ChannelCreate)
	htp.Register("/api/channels/{uuid}", http.MethodGet, handler.ChannelRead)
	htp.Register("/api/channels/{uuid}", http.MethodPut, handler.ChannelUpdate)
	htp.Register("/api/channels/{uuid}", http.MethodDelete, handler.ChannelDelete)

	htp.Register("/api/channels/{uuid}/messages", http.MethodGet, handler.ChannelMessagesRead)
	htp.Register("/api/channels/{uuid}/messages", http.MethodDelete, handler.ChannelMessagesDelete)

	htp.Register("/api/roles/@me", http.MethodGet, handler.RolesMe)
	htp.Register("/api/roles/create", http.MethodPost, handler.RolesCreate)
	htp.Register("/api/roles/{uuid}", http.MethodPut, handler.RoleUpdate)
	htp.Register("/api/roles/{uuid}", http.MethodDelete, handler.RoleDelete)

	htp.Register("/api/invites/@me", http.MethodGet, handler.InvitesMe)
	htp.Register("/api/invites/create", http.MethodPost, handler.InvitesCreate)
	htp.Register("/api/invites/{uuid}", http.MethodPut, handler.InviteUpdate)
	htp.Register("/api/invites/{uuid}", http.MethodDelete, handler.InviteDelete)

	htp.Register("/api/admin/audits.csv", http.MethodGet, handler.AuditsCsv)

	//
	// start server

	etc.StartServer()
}
