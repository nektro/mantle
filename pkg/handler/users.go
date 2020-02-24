package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"
)

func UsersMe(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
		"me":    user,
		"perms": ws.UserPerms{}.From(user),
	})
}

func UsersRead(w http.ResponseWriter, r *http.Request) {
	_, _, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	uu := r.URL.Path[len("/api/users/"):]
	u, ok := db.QueryUserByUUID(uu)
	writeAPIResponse(r, w, ok, http.StatusOK, u)
}

func UsersOnline(w http.ResponseWriter, r *http.Request) {
	_, _, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	writeAPIResponse(r, w, true, http.StatusOK, ws.AllOnlineIDs())
}
