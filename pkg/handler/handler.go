package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
)

func SaveOAuth2InfoCb(w http.ResponseWriter, r *http.Request, provider string, id string, name string, oa2resp map[string]interface{}) {
	ru := db.QueryUserBySnowflake(provider, id, name)
	util.Log("[user-login]", provider, id, ru.UUID, name)
	sess := etc.GetSession(r)
	sess.Values["user"] = ru.UUID
	sess.Save(r, w)
	db.QueryAssertUserName(ru.UUID, name)
}
