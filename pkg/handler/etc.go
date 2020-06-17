package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/idata"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go.etc/htp"
)

// EtcBadgeVersion is the handler for /api/etc/badges/version.svg
func EtcBadgeVersion(w http.ResponseWriter, r *http.Request) {
	hBadge(
		w, r,
		idata.Name,
		idata.Version,
		"blue",
	)
}

// EtcBadgeMembersOnline is the handler for /api/etc/badges/members_online.svg
func EtcBadgeMembersOnline(w http.ResponseWriter, r *http.Request) {
	hBadge(
		w, r,
		r.Host,
		strconv.FormatInt(ws.OnlineUserCount(), 10)+" online",
		"brightgreen",
	)
}

// EtcBadgeMembersTotal is the handler for /api/etc/badges/members_total.svg
func EtcBadgeMembersTotal(w http.ResponseWriter, r *http.Request) {
	hBadge(
		w, r,
		r.Host,
		strconv.FormatInt(db.User{}.MemberCount(), 10)+" members",
		"brightgreen",
	)
}

// EtcRoleColorCSS is the handler for /api/etc/role_colors.css
func EtcRoleColorCSS(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	controls.GetMemberUser(c, r, w)
	w.Header().Add("content-type", "text/css")
	ar := db.Role{}.AllSorted()
	for i := len(ar) - 1; i >= 0; i-- {
		item := ar[i]
		if len(item.Color) > 0 {
			fmt.Fprintln(w, `[data-role="`+item.UUID+`"] { color: `+item.Color+` !important; } /* `+item.Name+` */`)
		}
	}
	fmt.Fprintln(w)
	for i := len(ar) - 1; i >= 0; i-- {
		item := ar[i]
		if len(item.Color) > 0 {
			fmt.Fprintln(w, `[data-role="`+item.UUID+`"].bg { background-color: `+item.Color+` !important; } /* `+item.Name+` */`)
		}
	}
	fmt.Fprintln(w)
	for i := len(ar) - 1; i >= 0; i-- {
		item := ar[i]
		if len(item.Color) > 0 {
			fmt.Fprintln(w, `[data-role="`+item.UUID+`"].bg-bf::before { background-color: `+item.Color+` !important; } /* `+item.Name+` */`)
		}
	}
}
