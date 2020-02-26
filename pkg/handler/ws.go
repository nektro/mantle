package handler

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/valyala/fastjson"
)

// Websocket is the handler for /ws
func Websocket(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	wuser := ws.Connect(user, w, r)

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
			wuser.SendMessageRaw(map[string]string{
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
