package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	etc "github.com/nektro/go.etc"
	uuid "github.com/satori/go.uuid"

	. "github.com/nektro/go-util/alias"
)

func helperSaveCallbackInfo(w http.ResponseWriter, r *http.Request, provider string, id string, name string, oa2resp map[string]interface{}) {
	ru := queryUserBySnowflake(provider, id, name)
	log.Println("[user-login]", provider, id, ru.UUID, name)
	sess := etc.GetSession(r)
	sess.Values["user"] = ru.UUID
	sess.Save(r, w)
	queryAssertUserName(ru.UUID, name)
}

func newUUID() string {
	return strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
}

func createChannel(name string) string {
	id := etc.Database.QueryNextID(cTableChannels)
	uid := newUUID()
	log.Println("[channel-create]", uid, "#"+name)
	etc.Database.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%d', ?, '')", cTableChannels, id, uid, id), name)
	assertChannelMessagesTableExists(uid)
	return uid
}

func assertChannelMessagesTableExists(uid string) {
	etc.Database.CreateTable(F("%s%s", cTableMessagesPrefix, strings.Replace(uid, "-", "_", -1)), []string{"id", "int primary key"}, [][]string{
		{"uuid", "text"},
		{"sent_at", "text"},
		{"sent_by", "text"},
		{"text", "text"},
		{"test", "text"},
	})
}

func apiBootstrapRequireLogin(r *http.Request, w http.ResponseWriter, method string, assertMembership bool) (*sessions.Session, *RowUser, error) {
	if r.Method != method {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusMethodNotAllowed, "This action requires using HTTP "+method)
	}
	if method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return nil, nil, writeAPIResponse(r, w, false, http.StatusBadRequest, "Error parsing form data")
		}
	}

	sess := etc.GetSession(r)
	sessID := sess.Values["user"]

	if sessID == nil {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusUnauthorized, "Must login to access this resource")
	}

	userID := sessID.(string)
	user, _ := queryUserByUUID(userID)

	if assertMembership && !user.IsMember {
		return nil, nil, writeAPIResponse(r, w, false, http.StatusForbidden, "This action requires being a member of this server. ("+userID+")")
	}

	return sess, &user, nil
}

func writeAPIResponse(r *http.Request, w http.ResponseWriter, good bool, status int, message interface{}) error {
	resp := APIResponse{good, message}
	dat, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Header().Add("content-type", "application/json")
	w.Write(dat)
	return nil
}

func createRole(name string) string {
	id := etc.Database.QueryNextID(cTableRoles)
	uid := newUUID()
	log.Println("[role-create]", uid, name)
	etc.Database.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%d', ?, '', 1, 1)", cTableRoles, id, uid, id), name)
	return uid
}

func calculateUserPermissions(user *RowUser) *UserPerms {
	perms := UserPerms{}
	for _, item := range strings.Split(user.Roles, ",") {
		if item == "" {
			continue
		}
		role := roleCache[item]

		switch role.PermManageChannels {
		case PermDeny, PermAllow:
			perms.ManageChannels = GetPermColumnRealVal(role.PermManageChannels)
		}
		switch role.PermManageRoles {
		case PermDeny, PermAllow:
			perms.ManageRoles = GetPermColumnRealVal(role.PermManageRoles)
		}
	}
	return &perms
}
