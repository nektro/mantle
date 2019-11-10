package main

import (
	"log"
	"net/http"
	"strings"

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

func assertChannelMessagesTableExists(uid string) {
	etc.Database.CreateTable(F("%s%s", cTableMessagesPrefix, strings.Replace(uid, "-", "_", -1)), []string{"id", "int primary key"}, [][]string{
		{"uuid", "text"},
		{"sent_at", "text"},
		{"sent_by", "text"},
		{"text", "text"},
		{"test", "text"},
	})
}
