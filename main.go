package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/iconst"
	"github.com/nektro/mantle/pkg/idata"
	"github.com/nektro/mantle/pkg/itypes"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	"github.com/spf13/pflag"
	"github.com/valyala/fastjson"

	. "github.com/nektro/go-util/alias"

	_ "github.com/nektro/mantle/statik"
)

func main() {
	util.Log("Welcome to " + iconst.Name + " " + iconst.Version + ".")

	//
	pflag.IntVar(&idata.Config.Port, "port", 8000, "The port to bind the web server to.")
	etc.AppID = "mantle"
	etc.PreInit()

	//
	etc.Init("mantle", &idata.Config, "./invite", helperSaveCallbackInfo)

	//
	// database initialization

	db.Init()

	// for loop create channel message tables
	_chans := (db.Channel{}.All())
	for _, item := range _chans {
		assertChannelMessagesTableExists(item.UUID)
	}

	//
	// add default channel, if none exist

	if len(_chans) == 0 {
		createChannel("chat")
	}

	//
	// initialize server properties

	db.Props.SetDefault("name", iconst.Name)
	db.Props.SetDefault("owner", "")
	db.Props.SetDefault("public", "true")
	db.Props.Init()

	//
	// create server 'Owner' Role
	//		uneditable, and has all perms always

	pa := uint8(itypes.PermAllow)
	ws.RoleCache["o"] = db.Role{
		0, "o", 0, "Owner", "", pa, pa,
	}

	//
	// load roles into local cache

	for _, item := range (db.Role{}.All()) {
		ws.RoleCache[item.UUID] = item
	}

	//
	// setup graceful stop

	util.RunOnClose(func() {
		util.Log("Gracefully shutting down...")

		util.Log("Saving database to disk")
		db.DB.Close()

		util.Log("Closing all remaining active WebSocket connections")
		for _, item := range ws.ConnCache {
			item.Conn.Close()
		}

		util.Log("Done")
		os.Exit(0)
	})

	//
	// create http service

	http.HandleFunc("/invite", func(w http.ResponseWriter, r *http.Request) {
		_, user, _ := apiBootstrapRequireLogin(r, w, http.MethodGet, false)

		if db.Props.Get("public") == "true" {
			if user.IsMember == false {
				db.DB.Build().Up(iconst.TableUsers, "is_member", "1").Wh("uuid", user.UUID).Exe()
				util.Log("[user-join]", F("User %s just became a member and joined the server", user.UUID))
			}
			w.Header().Add("Location", "./chat/")
			w.WriteHeader(http.StatusFound)
		}
	})

	http.HandleFunc("/api/about", func(w http.ResponseWriter, r *http.Request) {
		dat, _ := json.Marshal(db.Props.GetAll())
		fmt.Fprint(w, string(dat))
	})

	http.HandleFunc("/api/users/@me", func(w http.ResponseWriter, r *http.Request) {
		_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
		if err != nil {
			return
		}
		writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
			"me":    user,
			"perms": ws.UserPerms{}.From(user),
		})
	})

	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		_, _, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
		if err != nil {
			return
		}
		uu := r.URL.Path[len("/api/users/"):]
		u, ok := queryUserByUUID(uu)
		writeAPIResponse(r, w, ok, http.StatusOK, u)
	})

	http.HandleFunc("/api/users/online", func(w http.ResponseWriter, r *http.Request) {
		_, _, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
		if err != nil {
			return
		}
		writeAPIResponse(r, w, true, http.StatusOK, listToArray(ws.Connected))
	})

	http.HandleFunc("/api/channels/@me", func(w http.ResponseWriter, r *http.Request) {
		writeAPIResponse(r, w, true, http.StatusOK, db.Channel{}.All())
	})

	http.HandleFunc("/api/channels/create", func(w http.ResponseWriter, r *http.Request) {
		_, user, err := apiBootstrapRequireLogin(r, w, http.MethodPost, true)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		if etc.AssertPostFormValuesExist(r, "name") != nil {
			fmt.Fprintln(w, "missing post value")
			return
		}
		cv, ok := ws.ConnCache[user.UUID]
		if !ok {
			fmt.Fprintln(w, "unable to find user in ws connection cache")
			return
		}
		if !cv.Perms.ManageChannels {
			fmt.Fprintln(w, "user missing 'Perms.ManageChannels'")
			return
		}
		name := r.Form.Get("name")
		cuid := createChannel(name)
		ws.BroadcastMessage(map[string]string{
			"type": "new-channel",
			"uuid": cuid,
			"name": name,
		})
	})

	//
	// create websocket service

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
		if err != nil {
			return
		}
		wuser := ws.User{}.From(r, w, user)
		ws.ConnCache[user.UUID] = wuser

		wuser.Connect()

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
				wuser.Conn.WriteJSON(map[string]string{
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
	})

	//
	// start server

	etc.StartServer(idata.Config.Port)
}
