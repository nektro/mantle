package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/gorilla/mux"
	"github.com/nektro/go.etc/htp"
)

// UsersMe is handler for /api/users/@me
func UsersMe(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
		"me":    user,
		"perms": ws.UserPerms{}.From(user),
	})
}

// UsersRead is handler for /api/users/{uuid}
func UsersRead(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	controls.GetMemberUser(c, r, w)
	uu := mux.Vars(r)["uuid"]
	u, ok := db.QueryUserByUUID(uu)
	writeAPIResponse(r, w, ok, http.StatusOK, u)
}

// UsersOnline is handler for /api/users/online
func UsersOnline(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	controls.GetMemberUser(c, r, w)
	writeAPIResponse(r, w, true, http.StatusOK, ws.AllOnlineIDs())
}

// UserUpdate is handler for /api/users/{uuid}/update
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	uu := mux.Vars(r)["uuid"]
	u, ok := db.QueryUserByUUID(uu)
	c.Assert(ok, "404: unable to find user with that uuid")
	controls.AssertFormKeysExist(c, r, "p_name")

	successCb := func(us *db.User, pk, pv string) {
		db.CreateAudit(db.ActionUserUpdate, user, us.UUID, pk, pv)
		writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
			"user":  us,
			"key":   pk,
			"value": pv,
		})
		ws.BroadcastMessage(map[string]interface{}{
			"type":  "user-update",
			"user":  us,
			"key":   pk,
			"value": pv,
		})
	}

	n := r.Form.Get("p_name")
	v := r.Form.Get("p_value")
	up := ws.UserPerms{}.From(user)
	if n != "nickname" {
		c.Assert(len(v) > 0, "400: missing form value p_value")
	}
	switch n {
	case "nickname":
		if user.UUID != u.UUID {
			return
		}
		u.SetNickname(v)
		successCb(u, n, v)
	case "add_role":
		rl, ok := db.QueryRoleByUID(v)
		if !ok {
			return
		}
		c.Assert(up.ManageRoles, "403: users require the manage_roles permission to update roles")
		c.Assert(user.GetRolesSorted()[0].Position < rl.Position, "403: role rank must be higher to update")
		u.AddRole(v)
		successCb(u, n, v)
	case "remove_role":
		rl, ok := db.QueryRoleByUID(v)
		if !ok {
			return
		}
		c.Assert(up.ManageRoles, "403: users require the manage_roles permission to update roles")
		c.Assert(user.GetRolesSorted()[0].Position < rl.Position, "403: role rank must be higher to update")
		u.RemoveRole(v)
		successCb(u, n, v)
	}
}
