package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/sessions"
	"github.com/nektro/go-util/alias"
	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
)

var formMethods = []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

func apiBootstrapRequireLogin(r *http.Request, w http.ResponseWriter, method string, assertMembership bool) (*sessions.Session, *db.User, error) {
	if r.Method != method {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusMethodNotAllowed, "This action requires using HTTP "+method)
	}
	if util.Contains(formMethods, method) {
		r.Method = http.MethodPost
		err := r.ParseMultipartForm(1024 * 1024 * 10)
		if err != nil {
			return nil, nil, writeAPIResponse(r, w, false, http.StatusBadRequest, "Error parsing form data. "+err.Error())
		}
		r.Method = method
	}

	sess := etc.GetSession(r)
	sessID := sess.Values["user"]

	if sessID == nil {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusUnauthorized, "Must login to access this resource")
	}

	userID := sessID.(string)
	user, _ := db.QueryUserByUUID(userID)

	if assertMembership && !user.IsMember {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusForbidden, "This action requires being a member of this server.")
	}

	return sess, user, nil
}

func writeAPIResponse(r *http.Request, w http.ResponseWriter, good bool, status int, message interface{}) error {
	resp := map[string]interface{}{
		"success": good,
		"message": message,
	}
	dat, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")
	fmt.Fprintln(w, string(dat))
	if !good {
		return errors.New("")
	}
	return nil
}

func hGrabInt(s string) (string, int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	return s, n, err
}

func hBadge(w http.ResponseWriter, r *http.Request, l, m, c string) {
	l = strings.ReplaceAll(l, " ", "_")
	m = strings.ReplaceAll(m, " ", "_")
	w.Header().Add("location", alias.F("https://img.shields.io/badge/%s-%s-%s", l, m, c))
	w.WriteHeader(http.StatusFound)
}

func hGrabFormStrings(r *http.Request, w http.ResponseWriter, s ...string) error {
	for _, item := range s {
		if !(len(r.Form.Get(item)) > 0) {
			writeAPIResponse(r, w, false, http.StatusBadRequest, "missing form value '"+item+"'.")
			return alias.E("missing " + item + " in form")
		}
	}
	return nil
}

func uHighLow(a, b int) (int, int) {
	if a >= b {
		return a, b
	}
	return b, a
}
