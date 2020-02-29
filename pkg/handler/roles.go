package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/ws"
)

// RolesRead reads info about channel
func RolesRead(w http.ResponseWriter, r *http.Request) {
	_, _, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	writeAPIResponse(r, w, true, http.StatusOK, ws.RoleCache)
}
