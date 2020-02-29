package handler

import (
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"
)

func EtcBadgeMembersOnline(w http.ResponseWriter, r *http.Request) {
	hBadge(
		w, r,
		r.Host,
		strconv.FormatInt(ws.OnlineUserCount(), 10)+" online",
		"brightgreen",
	)
}

func EtcBadgeMembersTotal(w http.ResponseWriter, r *http.Request) {
	hBadge(
		w, r,
		r.Host,
		strconv.FormatInt(db.User{}.MemberCount(), 10)+" members",
		"brightgreen",
	)
}
