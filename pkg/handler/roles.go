package handler

import (
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/colors.v1"
)

// RolesMe reads info about channel
func RolesMe(w http.ResponseWriter, r *http.Request) {
	writeAPIResponse(r, w, true, http.StatusOK, db.Role{}.All())
}

// RolesCreate reads info about channel
func RolesCreate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPost, true)
	if err != nil {
		return
	}
	n := r.Form.Get("name")
	if !(len(n) > 0) {
		writeAPIResponse(r, w, false, http.StatusBadRequest, "missing form value 'p_name'.")
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageRoles {
		writeAPIResponse(r, w, false, http.StatusForbidden, "users require the manage_roles permission to update roles.")
		return
	}
	nr := db.CreateRole(n)
	w.WriteHeader(http.StatusCreated)
	ws.BroadcastMessage(map[string]interface{}{
		"type": "role-new",
		"role": nr,
	})
}

// RoleUpdate updates info about this role
func RoleUpdate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPut, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageRoles {
		return
	}
	if hGrabFormStrings(r, w, "p_name") != nil {
		return
	}
	uu := mux.Vars(r)["uuid"]
	rl, ok := db.QueryRoleByUID(uu)
	if !ok {
		return
	}

	successCb := func(rs *db.Role, pk, pv string) {
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
		pH, pL := uHighLow(rl.Position, i)
		allR := db.Role{}.AllSorted()
		for d, item := range allR {
			o := d + 1
			if o < pL {
				continue
			}
			if o > pH {
				continue
			}
			// role moving down
			if pL == rl.Position {
				if o == pL {
					continue
				}
				if o == pH {
					rl.SetPosition(i)
					continue
				}
				item.SetPosition(o - 1)
			}
			// role moving up
			if pL == i {
				if o == pH {
					rl.SetPosition(i)
					continue
				}
				item.SetPosition(o + 1)
			}
		}
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
	}
}
