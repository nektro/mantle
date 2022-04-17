package handler

import (
	"net/http"

	"github.com/nektro/go-util/util"
	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go.etc/dbt"
	"github.com/nektro/go.etc/htp"
	"github.com/valyala/fastjson"
)

// Websocket is the handler for /ws
func Websocket(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	wuser, err := ws.Connect(user, w, r)
	if err != nil {
		return
	}

	// message intake loop
	for {
		// Read message from browser
		_, msg, err := wuser.ReadMessage()
		if err != nil {
			break
		}

		// broadcast message to all connected clients
		smg, err := fastjson.ParseBytes(msg)
		if err != nil {
			continue
		}
		switch string(smg.GetStringBytes("type")) {
		case "ping":
			// do nothing, keep connection alive
			wuser.SendWsMessage(map[string]interface{}{
				"type": "pong",
			})

		case "message":
			inch := dbt.UUID(string(smg.GetStringBytes("in")))
			c, ok := db.QueryChannelByUUID(inch)
			if !ok {
				continue
			}
			wuser.SendMessage(c, string(smg.GetStringBytes("message")))
		case "voice-connect", "voice-disconnect", "voice-data":
			ws.BroadcastMessageRaw(smg)
		default:
			util.LogError("ws: unhandled event discarded:", string(smg.GetStringBytes("type")))
		}
	}

	wuser.Disconnect()
}
