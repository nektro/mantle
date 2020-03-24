package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
	w.WriteHeader(http.StatusCreated)
	ws.BroadcastMessage(map[string]interface{}{
		"type":    "channel-new",
		"channel": nch,
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

// ChannelMessagesRead reads message data from channel
func ChannelMessagesRead(w http.ResponseWriter, r *http.Request) {
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

// ChannelMessagesDelete reads message data from channel
func ChannelMessagesDelete(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodDelete, true)
	if err != nil {
		fmt.Fprintln(w, 1)
		return
	}
	c, ok := db.QueryChannelByUUID(mux.Vars(r)["uuid"])
	if !ok {
		fmt.Fprintln(w, 2)
		return
	}
	actioned := []string{}
	for _, item := range r.Form["ids"] {
		if !db.IsUID(item) {
			continue
		}
		user.DeleteMessage(c, item)
		actioned = append(actioned, item)
	}
	ws.BroadcastMessage(map[string]interface{}{
		"type":     "message-delete",
		"channel":  c.UUID,
		"affected": actioned,
	})
	writeAPIResponse(r, w, true, http.StatusOK, actioned)
}

// ChannelUpdate updates info about this channel
func ChannelUpdate(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPut, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageChannels {
		return
	}
	if hGrabFormStrings(r, w, "p_name") != nil {
		return
	}
	uu := mux.Vars(r)["uuid"]
	ch, ok := db.QueryChannelByUUID(uu)
	if !ok {
		return
	}

	successCb := func(rs *db.Channel, pk, pv string) {
		writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
			"channel": rs,
			"key":     pk,
			"value":   pv,
		})
		ws.BroadcastMessage(map[string]interface{}{
			"type":    "channel-update",
			"channel": rs,
			"key":     pk,
			"value":   pv,
		})
	}

	n := r.Form.Get("p_name")
	v := r.Form.Get("p_value")
	switch n {
	case "name":
		if len(v) == 0 {
			return
		}
		ch.SetName(v)
		successCb(ch, n, v)
	case "position":
		i, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		pH, pL := uHighLow(ch.Position, i)
		allR := db.Role{}.AllSorted()
		for d, item := range allR {
			o := d + 1
			if o < pL {
				continue
			}
			if o > pH {
				continue
			}
			// role moving down
			if pL == ch.Position {
				if o == pL {
					continue
				}
				if o == pH {
					ch.SetPosition(i)
					continue
				}
				item.SetPosition(o - 1)
			}
			// role moving up
			if pL == i {
				if o == pH {
					ch.SetPosition(i)
					continue
				}
				item.SetPosition(o + 1)
			}
		}
		successCb(ch, n, v)
	case "description":
		if len(v) == 0 {
			return
		}
		ch.SetDescription(v)
		successCb(ch, n, v)
	case "history_off":
		b, err := strconv.ParseBool(v)
		if err != nil {
			return
		}
		ch.EnableHistory(b)
		successCb(ch, n, v)
	}
}

// ChannelDelete updates info about this channel
func ChannelDelete(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodDelete, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ManageChannels {
		return
	}
	uu := mux.Vars(r)["uuid"]
	v, ok := db.QueryChannelByUUID(uu)
	if !ok {
		return
	}
	v.Delete()
	ws.BroadcastMessage(map[string]interface{}{
		"type":    "channel-delete",
		"channel": uu,
	})
}
