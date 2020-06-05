package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	"github.com/nektro/go.etc/htp"
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

// InviteGet is handler for /
func InviteGet(w http.ResponseWriter, r *http.Request) {
	etc.WriteHandlebarsFile(r, w, "/invite.hbs", map[string]interface{}{
		"data": db.Props.GetAll(),
		"code": r.URL.Query().Get("code"),
	})
}

// InvitePost is handler for /invite
func InvitePost(w http.ResponseWriter, r *http.Request) {
	if ok, _ := strconv.ParseBool(db.Props.Get("public")); !ok {
		s := etc.GetSession(r)
		s.Values["code"] = r.URL.Query().Get("code")
		s.Save(r, w)
	}
	w.Header().Add("Location", "./login")
	w.WriteHeader(http.StatusFound)
}

// Verify is handler for /verify
func Verify(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetUser(c, r)
	c.RedirectIf(user.IsMember, "./chat/")

	if o, _ := strconv.ParseBool(db.Props.Get("public")); o {
		if !user.IsMember {
			user.SetAsMember(true)
		}
		c.RedirectIf(true, "./chat/")
		return
	}

	code, ok := sess.Values["code"].(string)
	c.Assert(ok, "400: invite code required to enter")

	inv, ok := db.QueryInviteByCode(code)
	c.Assert(ok, "400: invalid invite code")
	c.Assert(!inv.IsFrozen, "400: invite is frozen and can not be used")
	c.Assert(!(inv.MaxUses > 0 && inv.Uses >= inv.MaxUses), "400: invite use count has been exceeded")

	switch inv.Mode {
	case 0:
		// permanent
	case 1:
		//
	case 2:
		s := time.Since(inv.ExpiresOn.T())
		if s > 0 {
			writeAPIResponse(r, w, false, http.StatusBadRequest, "invite is expired")
			return
		}
	}

	inv.Use()
	user.SetAsMember(true)
	for _, item := range inv.GivenRoles {
		user.AddRole(item)
	}
	c.RedirectIf(true, "./chat/")
}

func Chat(w http.ResponseWriter, r *http.Request) {
	etc.WriteHandlebarsFile(r, w, "/chat/index.hbs", map[string]interface{}{
		"data": db.Props.GetAll(),
	})
}

// ApiAbout is handler for /api/about
func ApiAbout(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("all")) > 0 {
		c := htp.GetController(r)
		u := controls.GetMemberUser(c, r)
		c.Assert(u.HasRole("o"), "403: resource requires Authorization or to be server owner to access")
		writeAPIResponse(r, w, true, http.StatusOK, db.Props.GetAll())
		return
	}
	writeAPIResponse(r, w, true, http.StatusOK, db.Props.GetSome("name", "owner", "public", "description", "cover_photo", "profile_photo", "version"))
}

func ApiPropertyUpdate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r)
	controls.AssertFormKeysExist(c, r, "p_name", "p_value")

	n := r.Form.Get("p_name")
	v := r.Form.Get("p_value")
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageServer, "403: users require the manage_server permission to update properties")
	c.Assert(db.Props.Has(n), "400: specified property does not exist")

	db.Props.Set(n, v)
	db.CreateAudit(db.ActionSettingUpdate, user, "", n, v)
	writeAPIResponse(r, w, true, http.StatusOK, []string{n, v})
}
