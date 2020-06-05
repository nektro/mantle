package controls

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nektro/mantle/pkg/idata"

	"github.com/dgrijalva/jwt-go"
	"github.com/nektro/go.etc/htp"
)

// GetJWTClaims reads the Bearer token from the Request and asserts it is valid
func GetJWTClaims(c *htp.Controller, r *http.Request) map[string]interface{} {
	bearer := ""
	for _, item := range []func(*http.Request) string{tokenFromHeader, tokenFromCookie, tokenFromQuery} {
		bearer = item(r)
		if len(bearer) > 0 {
			break
		}
	}
	c.Assert(len(bearer) > 0, "400: bearer token missing")

	token, err := jwt.Parse(bearer, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method: " + fmt.Sprintf("%v", t.Header["alg"]))
		}
		return []byte(idata.Config.JWTSecret), nil
	})
	c.Assert(err == nil, "401: jwt error: "+fmt.Sprintf("%v", err))
	c.Assert(token.Valid, "401: jwt invalid")

	claims, ok := token.Claims.(jwt.MapClaims)
	c.Assert(ok, "401: jwt claims cast error")
	c.Assert(claims.Valid() == nil, "401: jwt claims invalid")

	return claims
}

func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) <= 7 || bearer[0:6] != "Bearer" {
		return ""
	}
	return bearer[7:]
}

func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func tokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("jwt")
}
