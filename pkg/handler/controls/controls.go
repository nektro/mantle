package controls

import (
	"net/http"

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
