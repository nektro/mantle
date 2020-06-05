package controls

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/sessions"
	"github.com/nektro/go-util/arrays/stringsu"
	etc "github.com/nektro/go.etc"
	"github.com/nektro/go.etc/htp"
)

// AssertFormKeysExist asserts Request.Form keys exist in htp
func AssertFormKeysExist(c *htp.Controller, r *http.Request, s ...string) map[string]string {
	res := map[string]string{}
	for _, item := range s {
		v := r.Form.Get(item)
		c.Assert(len(v) > 0, "400: missing post value: "+item)
		res[item] = v
	}
	return res
}

// GetSession grabs session
func GetSession(c *htp.Controller, r *http.Request) *sessions.Session {
	return etc.GetSession(r)
}

var formMethods = []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

// GetUser asserts a user is logged in
func GetUser(c *htp.Controller, r *http.Request) *db.User {
	s := GetSession(c, r)
	sessID := s.Values["user"]
	c.Assert(sessID != nil, "403: must login to access this resource")
	//
	userID := sessID.(string)
	user, _ := db.QueryUserByUUID(userID)

	method := r.Method
	if stringsu.Contains(formMethods, method) {
		r.Method = http.MethodPost
		r.ParseMultipartForm(0)
		r.Method = method
	}

	return user
}

// GetMemberUser asserts the user is a member and not banned
func GetMemberUser(c *htp.Controller, r *http.Request) *db.User {
	u := GetUser(c, r)
	c.Assert(u.IsMember, "403: you are not a member of this server")
	c.Assert(!u.IsBanned, "403: you are banned")
	return u
}
