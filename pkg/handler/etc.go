package handler

import (
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"
)

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
