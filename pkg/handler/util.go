package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/sessions"
	"github.com/nektro/go-util/arrays/stringsu"
	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	"github.com/nektro/go.etc/htp"
	sdrie "github.com/nektro/go.sdrie"

	. "github.com/nektro/go-util/alias"
)

var (
	formMethods = []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	badgeCache  = sdrie.New()
)

// Init sets up this package
func Init() {
	htp.ErrorHandleFunc = func(w http.ResponseWriter, r *http.Request, data string) {
		code, _ := strconv.ParseInt(data[:3], 10, 32)
		good := !(code >= 400)
		writeAPIResponse(r, w, good, int(code), data[5:])
	}
}

func apiBootstrapRequireLogin(r *http.Request, w http.ResponseWriter, method string, assertMembership bool) (*sessions.Session, *db.User, error) {
	if r.Method != method {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusMethodNotAllowed, "This action requires using HTTP "+method)
	}
	if stringsu.Contains(formMethods, method) {
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
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	dat, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(dat))
	if !good {
		return E(F("%v", message))
	}
	return nil
}

func hGrabInt(s string) (string, int64, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	return s, n, err
}

func hBadge(w http.ResponseWriter, r *http.Request, l, m, c string) {
	l = strings.ReplaceAll(l, " ", "_")
	m = strings.ReplaceAll(m, " ", "_")
	k := l + "-" + m + "-" + c
	w.Header().Add("content-type", "image/svg+xml")
	if badgeCache.Has(k) {
		fmt.Fprintln(w, string(badgeCache.Get(k).([]byte)))
		return
	}
	req, _ := http.NewRequest(http.MethodGet, "https://img.shields.io/badge/"+k, nil)
	bys, err := util.DoHttpFetch(req)
	if err != nil {
		return
	}
	badgeCache.Set(k, bys, time.Minute*10)
	fmt.Fprintln(w, string(bys))
}

func hBetween(x, l, h int) bool {
	if x < l {
		return false
	}
	if x > h {
		return false
	}
	return true
}
