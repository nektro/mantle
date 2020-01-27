package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nektro/mantle/pkg/idata"
	"github.com/nektro/mantle/pkg/itypes"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	"github.com/spf13/pflag"
	"github.com/valyala/fastjson"

	. "github.com/nektro/go-util/alias"

	_ "github.com/nektro/mantle/statik"
)


func main() {
	util.Log("Welcome to " + Name + ".")

	//
	pflag.IntVar(&idata.Config.Port, "port", 8000, "The port to bind the web server to.")
	etc.PreInit()

	//
	etc.Init("mantle", &idata.Config, "./invite", helperSaveCallbackInfo)

	//
	// database initialization

	etc.Database.CreateTableStruct(cTableSettings, itypes.Setting{})
	etc.Database.CreateTableStruct(cTableUsers, itypes.User{})
	etc.Database.CreateTableStruct(cTableChannels, itypes.Channel{})
	etc.Database.CreateTableStruct(cTableRoles, itypes.Role{})
	etc.Database.CreateTableStruct(cTableChannelRolePerms, itypes.ChannelRolePerms{})

	// for loop create channel message tables
	_chans := queryAllChannels()
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

	props.SetDefault("name", Name)
	props.SetDefault("owner", "")
	props.SetDefault("public", "true")
	props.Init()

	//
	// create server 'Owner' Role
	//		uneditable, and has all perms always

	pa := uint8(PermAllow)
	idata.RoleCache["o"] = itypes.Role{
		0, "o", 0, "Owner", "", pa, pa,
	}

	//
	// load roles into local cache

	for _, item := range queryAllRoles() {
		idata.RoleCache[item.UUID] = item
	}

	//
	// setup graceful stop

	util.RunOnClose(func() {
		util.Log("Gracefully shutting down...")

		util.Log("Saving database to disk")
		etc.Database.Close()

		util.Log("Closing all remaining active WebSocket connections")
		for _, item := range idata.WsConnCache {
			item.Conn.Close()
		}

		util.Log("Done")
		os.Exit(0)
	})

	//
	// create http service

	http.HandleFunc("/invite", func(w http.ResponseWriter, r *http.Request) {
		_, user, _ := apiBootstrapRequireLogin(r, w, http.MethodGet, false)

		if props.Get("public") == "true" {
			if user.IsMember == false {
				etc.Database.Build().Up(cTableUsers, "is_member", "1").Wh("uuid", user.UUID).Exe()
				util.Log("[user-join]", F("User %s just became a member and joined the server", user.UUID))
			}
			w.Header().Add("Location", "./chat/")
			w.WriteHeader(http.StatusFound)
		}
	})

	http.HandleFunc("/api/about", func(w http.ResponseWriter, r *http.Request) {
		dat, _ := json.Marshal(props.GetAll())
		fmt.Fprint(w, string(dat))
	})

	http.HandleFunc("/api/users/@me", func(w http.ResponseWriter, r *http.Request) {
		_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
		if err != nil {
			return
		}
		writeAPIResponse(r, w, true, http.StatusOK, map[string]interface{}{
			"me":    user,
			"perms": calculateUserPermissions(user),
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
		writeAPIResponse(r, w, true, http.StatusOK, listToArray(idata.Connected))
	})

	http.HandleFunc("/api/channels/@me", func(w http.ResponseWriter, r *http.Request) {
		channels := queryAllChannels()
		writeAPIResponse(r, w, true, http.StatusOK, channels)
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
		cv, ok := idata.WsConnCache[user.UUID]
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
		broadcastMessage(map[string]string{
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
		conn, _ := idata.WsUpgrader.Upgrade(w, r, nil)
		perms := calculateUserPermissions(user)
		idata.WsConnCache[user.UUID] = itypes.ConnCacheValue{conn, user, perms}

		// connect
		if !listHas(idata.Connected, user.UUID) {
			idata.Connected.PushBack(user.UUID)
			broadcastMessage(map[string]string{
				"type": "user-connect",
				"user": user.UUID,
			})
		}
		// message intake loop
		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
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
				conn.WriteJSON(map[string]string{
					"type": "pong",
				})
			case "message":
				broadcastMessage(map[string]string{
					"type":    "message",
					"in":      string(smg.GetStringBytes("in")),
					"from":    user.UUID,
					"message": string(smg.GetStringBytes("message")),
				})
			}
		}
		// disconnect
		if listHas(idata.Connected, user.UUID) {
			delete(idata.WsConnCache, user.UUID)
			listRemove(idata.Connected, user.UUID)
			broadcastMessage(map[string]string{
				"type": "user-disconnect",
				"user": user.UUID,
			})
		}
	})

	//
	// start server

	etc.StartServer(idata.Config.Port)
}
