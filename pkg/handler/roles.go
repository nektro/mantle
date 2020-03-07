package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"
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
