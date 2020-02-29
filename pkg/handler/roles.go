package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
)

// RolesRead reads info about channel
func RolesRead(w http.ResponseWriter, r *http.Request) {
	writeAPIResponse(r, w, true, http.StatusOK, db.Role{}.All())
}
