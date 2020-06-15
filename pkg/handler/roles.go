package handler

import (
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/go-playground/colors"
	"github.com/gorilla/mux"
	"github.com/nektro/go.etc/htp"
)

// RolesMe is the handler for /api/roles/@me
func RolesMe(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	controls.GetMemberUser(c, r, w)
	writeAPIResponse(r, w, true, http.StatusOK, db.Role{}.All())
}

// RolesCreate reads info about channel
func RolesCreate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	controls.AssertFormKeysExist(c, r, "name")
	n := r.Form.Get("name")

	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageRoles, "403: users require the manage_roles permission to update roles")

	nr := db.CreateRole(n)
	db.CreateAudit(db.ActionRoleCreate, user, nr.UUID, "", "")
	w.WriteHeader(http.StatusCreated)
	ws.BroadcastMessage(map[string]interface{}{
		"type": "role-new",
		"role": nr,
	})
}

// RoleUpdate updates info about this role
func RoleUpdate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageRoles, "403: users require the manage_roles permission to update roles")
	controls.AssertFormKeysExist(c, r, "p_name")

	uu := mux.Vars(r)["uuid"]
	rl, ok := db.QueryRoleByUID(uu)
	c.Assert(ok, "404: unable to find role with that uuid")
	c.Assert(user.GetRolesSorted()[0].Position < rl.Position, "403: role rank must be higher to update")

	successCb := func(rs *db.Role, pk, pv string) {
		db.CreateAudit(db.ActionRoleUpdate, user, rs.UUID, pk, pv)
		writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
			"role":  rs,
			"key":   pk,
			"value": pv,
		})
		ws.BroadcastMessage(map[string]interface{}{
			"type":  "role-update",
			"role":  rs,
			"key":   pk,
			"value": pv,
		})
	}
	processPerm := func(n, v string, rs *db.Role, f func(int)) {
		_, a, err := hGrabInt(v)
		if err != nil {
			return
		}
		b := int(a)
		if !hBetween(b, 0, 2) {
			return
		}
		f(b)
		successCb(rs, n, v)
	}

	n := r.Form.Get("p_name")
	v := r.Form.Get("p_value")
	switch n {
	case "name":
		if len(v) == 0 {
			return
		}
		rl.SetName(v)
		successCb(rl, n, v)
	case "color":
		_, err := colors.Parse(v)
		if err != nil {
			return
		}
		rl.SetColor(v)
		successCb(rl, n, v)
	case "position":
		i, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		rl.MoveTo(i)
		successCb(rl, n, v)
	case "distinguish":
		b, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		rl.SetDistinguish(b)
		successCb(rl, n, v)
	case "perm_manage_server":
		processPerm(n, v, rl, func(x int) {
			rl.SetPermMngServer(x)
		})
	case "perm_manage_channels":
		processPerm(n, v, rl, func(x int) {
			rl.SetPermMngChannels(x)
		})
	case "perm_manage_roles":
		processPerm(n, v, rl, func(x int) {
			rl.SetPermMngRoles(x)
		})
	case "perm_manage_invites":
		processPerm(n, v, rl, func(x int) {
			rl.SetPermMngInvites(x)
		})
	case "perm_view_audits":
		processPerm(n, v, rl, func(x int) {
			rl.SetPermViewAudits(x)
		})
	}
}

// RoleDelete updates info about this role
func RoleDelete(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageRoles, "403: users require the manage_roles permission to update roles")

	uu := mux.Vars(r)["uuid"]
	rl, ok := db.QueryRoleByUID(uu)
	c.Assert(ok, "404: unable to find role with that uuid")
	c.Assert(user.GetRolesSorted()[0].Position < rl.Position, "403: role rank must be higher to update")

	us := rl.Delete()
	db.CreateAudit(db.ActionRoleDelete, user, rl.UUID, "", "")
	for _, item := range us {
		ws.BroadcastMessage(map[string]interface{}{
			"type":  "user-update",
			"user":  item,
			"key":   "remove_role",
			"value": rl.UUID,
		})
	}
	ws.BroadcastMessage(map[string]interface{}{
		"type": "role-delete",
		"role": uu,
	})
	w.WriteHeader(http.StatusNoContent)
}
