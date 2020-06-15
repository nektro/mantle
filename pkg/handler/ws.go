package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

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
		_, msg, err := wuser.Conn.ReadMessage()
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
			wuser.SendMessageRaw(map[string]interface{}{
				"type": "pong",
			})

		case "message":
			c, ok := db.QueryChannelByUUID(string(smg.GetStringBytes("in")))
			if !ok {
				continue
			}
			wuser.SendMessage(c, string(smg.GetStringBytes("message")))
		}
	}

	wuser.Disconnect()
}
