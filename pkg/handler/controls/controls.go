package controls

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"

	"github.com/nektro/go-util/arrays/stringsu"
	etc "github.com/nektro/go.etc"
	"github.com/nektro/go.etc/htp"
	"github.com/nektro/go.etc/jwt"
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

var formMethods = []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

// GetUser asserts a user is logged in
func GetUser(c *htp.Controller, r *http.Request, w http.ResponseWriter) *db.User {
	l := GetJWTClaims(c, r, w)
	//
	userID := l["sub"].(string)
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
func GetMemberUser(c *htp.Controller, r *http.Request, w http.ResponseWriter) *db.User {
	u := GetUser(c, r, w)
	c.Assert(u.IsMember, "403: you are not a member of this server")
	c.Assert(!u.IsBanned, "403: you are banned")
	return u
}

// GetJWTClaims reads the Bearer token from the Request and asserts it is valid
func GetJWTClaims(c *htp.Controller, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	claims, ok := jwt.VerifyRequest(r, etc.JWTSecret)
	c.Assert(ok, "401: jwt missing/invalid")
	w.Header().Add("x-iss", claims["iss"].(string))
	w.Header().Add("x-sub", claims["sub"].(string))
	return claims
}
