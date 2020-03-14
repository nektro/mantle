package handler

import (
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/colors.v1"
)

// RolesRead reads info about channel
func RolesRead(w http.ResponseWriter, r *http.Request) {
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
		writeAPIResponse(r, w, false, http.StatusForbidden, "users require the manage_server permission to update properties.")
		return
	}
	nr := db.CreateRole(n)
	ws.BroadcastMessage(map[string]interface{}{
		"type": "new-role",
		"role": nr,
	})
}

// RoleUpdate reads info about channel
func RoleUpdate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPut, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageRoles {
		writeAPIResponse(r, w, false, http.StatusForbidden, "users require the manage_server permission to update properties.")
		return
	}
	if hGrabFormStrings(r, w, "p_name") != nil {
		return
	}
	uu := mux.Vars(r)["uuid"]
	rl, ok := db.QueryRoleByUID(uu)
	if !ok {
		writeAPIResponse(r, w, false, http.StatusBadRequest, "missing uuid url parameter")
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
			writeAPIResponse(r, w, false, http.StatusBadRequest, "error parsing color: "+err.Error())
			return
		}
		rl.SetColor(v)
		successCb(rl, n, v)
	case "position":
		i, err := strconv.Atoi(v)
		if err != nil {
			writeAPIResponse(r, w, false, http.StatusBadRequest, "error parsing position: "+err.Error())
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
	default:
		writeAPIResponse(r, w, false, http.StatusBadRequest, "invalid p_name parameter")
	}
}
