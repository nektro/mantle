package handler

import (
	"net/http"
	"strconv"

	"github.com/nektro/go.etc/dbt"
	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go.etc/htp"
)

// ChannelsMe is the handler for /api/channels/@me
func ChannelsMe(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	controls.GetMemberUser(c, r, w)
	writeAPIResponse(r, w, true, http.StatusOK, db.Channel{}.All())
}

// ChannelCreate is the handler for /api/channels/create
func ChannelCreate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	controls.AssertFormKeysExist(c, r, "name")

	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageChannels, "403: action requires the manage_channels permission")

	name := r.Form.Get("name")
	nch := db.CreateChannel(name)
	db.CreateAudit(db.ActionChannelCreate, user, nch.UUID, "", "")
	w.WriteHeader(http.StatusCreated)
	ws.BroadcastMessage(map[string]interface{}{
		"type":    "channel-new",
		"channel": nch,
	})
}

// ChannelRead reads info about channel
func ChannelRead(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	controls.GetMemberUser(c, r, w)
	uu := controls.GetUIDFromPath(c, r)
	ch, ok := db.QueryChannelByUUID(uu)
	writeAPIResponse(r, w, ok, http.StatusOK, ch)
}

// ChannelMessagesRead reads message data from channel
func ChannelMessagesRead(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	controls.GetMemberUser(c, r, w)
	ch, ok := db.QueryChannelByUUID(controls.GetUIDFromPath(c, r))
	c.Assert(ok, "404: unable to find channel with this uuid")

	slm, lmn, err := hGrabInt(r.URL.Query().Get("limit"))
	c.Assert(len(slm) == 0 || err == nil, "400: error parsing limit query parameter")
	if lmn == 0 {
		lmn = 50
	}
	c.Assert(lmn > 0, "400: limit minimum is 1")
	c.Assert(lmn <= 50, "400: limit max is 50")

	aft := dbt.UUID(r.URL.Query().Get("after"))
	msgs := ch.QueryMsgAfterUID(aft, int(lmn))
	writeAPIResponse(r, w, true, http.StatusOK, msgs)
}

// ChannelMessagesDelete reads message data from channel
func ChannelMessagesDelete(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	ch, ok := db.QueryChannelByUUID(controls.GetUIDFromPath(c, r))
	c.Assert(ok, "404: unable to find channel with this uuid")

	actioned := []string{}
	for _, item := range r.Form["ids"] {
		if !dbt.IsUUID(dbt.UUID(item)) {
			continue
		}
		user.DeleteMessage(ch, dbt.UUID(item))
		actioned = append(actioned, item)
	}
	ws.BroadcastMessage(map[string]interface{}{
		"type":     "message-delete",
		"channel":  ch.UUID,
		"affected": actioned,
	})
	w.WriteHeader(http.StatusNoContent)
}

// ChannelUpdate updates info about this channel
func ChannelUpdate(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageChannels, "403: action requires the manage_channels permission")
	controls.AssertFormKeysExist(c, r, "p_name")

	uu := controls.GetUIDFromPath(c, r)
	ch, ok := db.QueryChannelByUUID(uu)
	c.Assert(ok, "404: unable to find channel with this uuid")

	successCb := func(rs *db.Channel, pk, pv string) {
		db.CreateAudit(db.ActionChannelUpdate, user, rs.UUID, pk, pv)
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
		ch.MoveTo(i)
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
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ManageChannels, "403: action requires the manage_channels permission")

	uu := controls.GetUIDFromPath(c, r)
	ch, ok := db.QueryChannelByUUID(uu)
	c.Assert(ok, "404: unable to find channel with this uuid")

	ch.Delete()
	db.CreateAudit(db.ActionChannelDelete, user, ch.UUID, "", "")
	ws.BroadcastMessage(map[string]interface{}{
		"type":    "channel-delete",
		"channel": uu,
	})
	w.WriteHeader(http.StatusNoContent)
}
