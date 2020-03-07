package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
)

// SaveOAuth2InfoCb saves info from go.oauth to user session cookie
func SaveOAuth2InfoCb(w http.ResponseWriter, r *http.Request, provider string, id string, name string, oa2resp map[string]interface{}) {
	ru := db.QueryUserBySnowflake(provider, id, name)
	util.Log("[user-login]", provider, id, ru.UUID, name)
	sess := etc.GetSession(r)
	sess.Values["user"] = ru.UUID
	sess.Save(r, w)
	ru.SetName(strings.ReplaceAll(name, " ", ""))
}

// InviteGet is handler for GET /invite
func InviteGet(w http.ResponseWriter, r *http.Request) {
	etc.WriteHandlebarsFile(r, w, "/invite.hbs", map[string]interface{}{
		"data": db.Props.GetAll(),
	})
}

// InvitePost is handler for POST /invite
func InvitePost(w http.ResponseWriter, r *http.Request) {
	if ok, _ := strconv.ParseBool(db.Props.Get("public")); ok {
		w.Header().Add("Location", "./login")
		w.WriteHeader(http.StatusFound)
		return
	}
}

// Verify is handler for /verify
func Verify(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, false)
	if err != nil {
		return
	}

	if o, _ := strconv.ParseBool(db.Props.Get("public")); o {
		if !user.IsMember {
			user.SetAsMember(true)
			util.Log("[user-join]", "User", user.UUID, "just became a member and joined the server")
		}
		w.Header().Add("Location", "./chat/")
		w.WriteHeader(http.StatusFound)
	}
}

// ApiAbout is handler for /api/about
func ApiAbout(w http.ResponseWriter, r *http.Request) {
	writeAPIResponse(r, w, true, http.StatusOK, db.Props.GetAll())
}

func ApiPropertyUpdate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPost, true)
	if err != nil {
		return
	}
	if hGrabFormStrings(r, w, "p_name", "p_value") != nil {
		return
	}
	n := r.Form.Get("p_name")
	v := r.Form.Get("p_value")
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageServer {
		writeAPIResponse(r, w, false, http.StatusForbidden, "users require the manage_server permission to update properties.")
		return
	}
	if !db.Props.Has(n) {
		writeAPIResponse(r, w, false, http.StatusBadRequest, "specified property does not exist.")
		return
	}
	db.Props.Set(n, v)
	writeAPIResponse(r, w, true, http.StatusOK, []string{n, v})
}
