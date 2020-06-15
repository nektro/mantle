package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/gorilla/mux"
	"github.com/nektro/go.etc/htp"
)

// InvitesMe reads info about channel
func InvitesMe(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageInvites {
		writeAPIResponse(r, w, true, http.StatusOK, []db.Invite{})
		return
	}
	writeAPIResponse(r, w, true, http.StatusOK, db.Invite{}.All())
}

// InvitesCreate reads info about channel
func InvitesCreate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageInvites, "403: users require the manage_invites permission to update invites")

	nr := db.CreateInvite()
	db.CreateAudit(db.ActionInviteCreate, user, nr.UUID, "", "")
	w.WriteHeader(http.StatusCreated)
	ws.BroadcastMessage(map[string]interface{}{
		"type":   "invite-new",
		"invite": nr,
	})
}

// InviteUpdate updates info about this invite
func InviteUpdate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageInvites, "403: users require the manage_invites permission to update invites")
	controls.AssertFormKeysExist(c, r, "p_name")

	uu := mux.Vars(r)["uuid"]
	iv, ok := db.QueryInviteByUID(uu)
	c.Assert(ok, "404: unable to find invite with that uuid")

	successCb := func(rs *db.Invite, pk, pv string) {
		db.CreateAudit(db.ActionInviteUpdate, user, rs.UUID, pk, pv)
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
		c.Assert(err == nil, "400: error parsing p_value")
		c.Assert(x >= 0, "400: p_value must be >= 0")
		iv.SetMaxUses(x)
		successCb(iv, n, v)
	}
}

// InviteDelete updates info about this invite
func InviteDelete(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageInvites, "403: users require the manage_invites permission to update invites")

	uu := mux.Vars(r)["uuid"]
	iv, ok := db.QueryInviteByUID(uu)
	c.Assert(ok, "404: unable to find invite with that uuid")

	iv.Delete()
	db.CreateAudit(db.ActionInviteDelete, user, iv.UUID, "", "")
	ws.BroadcastMessage(map[string]interface{}{
		"type":   "invite-delete",
		"invite": uu,
	})
	w.WriteHeader(http.StatusNoContent)
}
