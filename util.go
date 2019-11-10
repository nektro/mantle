package main

import (
	"log"
	"net/http"

	etc "github.com/nektro/go.etc"
)

func helperSaveCallbackInfo(w http.ResponseWriter, r *http.Request, provider string, id string, name string, oa2resp map[string]interface{}) {
	log.Println("[user-login]", provider, id, ru.UUID, name)
	sess := etc.GetSession(r)
	sess.Values["user"] = ru.UUID
	sess.Save(r, w)
}
