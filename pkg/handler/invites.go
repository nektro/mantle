package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/gorilla/mux"
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

// InvitesCreate reads info about channel
func InvitesCreate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPost, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageInvites {
		writeAPIResponse(r, w, false, http.StatusForbidden, "users require the manage_invites permission to update invites.")
		return
	}
	nr := db.CreateInvite()
	w.WriteHeader(http.StatusCreated)
	ws.BroadcastMessage(map[string]interface{}{
		"type":   "invite-new",
		"invite": nr,
	})
}

// InviteUpdate updates info about this invite
func InviteUpdate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPut, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageInvites {
		return
	}
	if hGrabFormStrings(r, w, "p_name") != nil {
		return
	}
	uu := mux.Vars(r)["uuid"]
	rl, ok := db.QueryInviteByUID(uu)
	if !ok {
		return
	}

	successCb := func(rs *db.Invite, pk, pv string) {
		writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
			"invite": rs,
			"key":    pk,
			"value":  pv,
		})
		ws.BroadcastMessage(map[string]interface{}{
			"type":   "invite-update",
			"invite": rs,
			"key":    pk,
			"value":  pv,
		})
	}

	n := r.Form.Get("p_name")
	v := r.Form.Get("p_value")
	switch n {
	case "max_uses":
		_, x, err := hGrabInt(v)
		if err != nil {
			return
		}
		rl.SetMaxUses(x)
		successCb(rl, n, v)
	}
}
