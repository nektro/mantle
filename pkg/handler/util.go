package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/itypes"

	"github.com/gorilla/sessions"
	etc "github.com/nektro/go.etc"
)

func apiBootstrapRequireLogin(r *http.Request, w http.ResponseWriter, method string, assertMembership bool) (*sessions.Session, *db.User, error) {
	if r.Method != method {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusMethodNotAllowed, "This action requires using HTTP "+method)
	}
	if method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return nil, nil, writeAPIResponse(r, w, false, http.StatusBadRequest, "Error parsing form data")
		}
	}

	sess := etc.GetSession(r)
	sessID := sess.Values["user"]

	if sessID == nil {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusUnauthorized, "Must login to access this resource")
	}

	userID := sessID.(string)
	user, _ := db.QueryUserByUUID(userID)

	if assertMembership && !user.IsMember {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusForbidden, "This action requires being a member of this server. ("+userID+")")
	}

	return sess, user, nil
}

func writeAPIResponse(r *http.Request, w http.ResponseWriter, good bool, status int, message interface{}) error {
	resp := itypes.APIResponse{good, message}
	dat, _ := json.Marshal(resp)
	w.Header().Add("content-type", "application/json")
	fmt.Fprintln(w, dat)
	w.WriteHeader(status)
	if !good {
		return errors.New("")
	}
	return nil
}
