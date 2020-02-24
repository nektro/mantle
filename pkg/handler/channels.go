package handler

import (
	"fmt"
	"net/http"

	etc "github.com/nektro/go.etc"
	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"
)

func ChannelsMe(w http.ResponseWriter, r *http.Request) {
	writeAPIResponse(r, w, true, http.StatusOK, db.Channel{}.All())
}

func ChannelCreate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPost, true)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	if etc.AssertPostFormValuesExist(r, "name") != nil {
		fmt.Fprintln(w, "missing post value")
		return
	}
	cv, ok := ws.UserCache[user.UUID]
	if !ok {
		fmt.Fprintln(w, "unable to find user in ws connection cache")
		return
	}
	if !cv.Perms.ManageChannels {
		fmt.Fprintln(w, "user missing 'Perms.ManageChannels'")
		return
	}
	name := r.Form.Get("name")
	cuid := db.CreateChannel(name)
	ws.BroadcastMessage(map[string]string{
		"type": "new-channel",
		"uuid": cuid,
		"name": name,
	})
}
