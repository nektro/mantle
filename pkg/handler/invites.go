package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"
)

// InvitesMe reads info about channel
func InvitesMe(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageInvites {
		return
	}
	writeAPIResponse(r, w, true, http.StatusOK, db.Invite{}.All())
}
