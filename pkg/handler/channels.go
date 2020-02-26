package handler

import (
	"fmt"
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/gorilla/mux"
	etc "github.com/nektro/go.etc"
)

// ChannelsMe is the handler for /api/channels/@me
func ChannelsMe(w http.ResponseWriter, r *http.Request) {
	writeAPIResponse(r, w, true, http.StatusOK, db.Channel{}.All())
}

// ChannelCreate is the handler for /api/channels/create
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
	nch := db.CreateChannel(name)
	ws.BroadcastMessage(map[string]interface{}{
		"type": "new-channel",
		"uuid": nch.UUID,
		"name": name,
	})
}

// ChannelRead reads info about channel
func ChannelRead(w http.ResponseWriter, r *http.Request) {
	_, _, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	uu := mux.Vars(r)["uuid"]
	u, ok := db.QueryChannelByUUID(uu)
	writeAPIResponse(r, w, ok, http.StatusOK, u)
}

// ChannelMessages reads message data from channel
func ChannelMessages(w http.ResponseWriter, r *http.Request) {
	_, _, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		fmt.Fprintln(w, 1)
		return
	}
	c, ok := db.QueryChannelByUUID(mux.Vars(r)["uuid"])
	if !ok {
		fmt.Fprintln(w, 2)
		return
	}
	_, lmn, err := hGrabInt(r.URL.Query().Get("limit"))
	if err != nil {
		lmn = 50
	}
	if lmn > 50 {
		fmt.Fprintln(w, 4)
		return
	}
	msgs := c.QueryMsgAfterUID(r.URL.Query().Get("after"), int(lmn))
	writeAPIResponse(r, w, true, http.StatusOK, msgs)
}
