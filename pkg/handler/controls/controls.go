package controls

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/sessions"
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

// GetUser asserts a user is logged in
func GetUser(c *htp.Controller, s *sessions.Session) *db.User {
	sessID := s.Values["user"]
	c.Assert(sessID != nil, "403: must login to access this resource")
	//
	userID := sessID.(string)
	user, _ := db.QueryUserByUUID(userID)
	return user
}

// AssertUserIsMember asserts the user is a member and not banned
func AssertUserIsMember(c *htp.Controller, u *db.User) {
	c.Assert(u.IsMember, "403: you are not a member of this server")
	c.Assert(!u.IsBanned, "403: you are banned")
}
