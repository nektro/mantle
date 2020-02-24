package handler

import (
	"net/http"
	"time"

	"github.com/nektro/mantle/pkg/ws"
	"github.com/valyala/fastjson"
)

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
			wuser.SendMessage(map[string]string{
				"type": "pong",
			})
		case "message":
			ws.BroadcastMessage(map[string]string{
				"type":    "message",
				"in":      string(smg.GetStringBytes("in")),
				"from":    user.UUID,
				"message": string(smg.GetStringBytes("message")),
				"at":      time.Now().UTC().Format("2 Jan 2006 15:04:05 MST"),
			})
		}
	}

	wuser.Disconnect()
}
