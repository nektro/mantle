package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nektro/mantle/pkg/itypes"

	"github.com/gorilla/websocket"
	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	oauth2 "github.com/nektro/go.oauth2"
	"github.com/spf13/pflag"

	. "github.com/nektro/go-util/alias"

	_ "github.com/nektro/mantle/statik"
)

var (
	config      *Config
	wsUpgrader  = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	wsConnCache = map[string]itypes.ConnCacheValue{}
	roleCache   = map[string]itypes.RowRole{}
	connected   = list.New()
)

type Config struct {
	Version   int               `json:"version"`
	Port      int               `json:"port"`
	Clients   []oauth2.AppConf  `json:"clients"`
	Providers []oauth2.Provider `json:"providers"`
}

func main() {
	util.Log("Welcome to " + Name + ".")

	//
	flagPort := pflag.Int("port", 0, "The port to bind the web server to.")
	etc.PreInit()

	//
	etc.Init("mantle", &config, "./invite", helperSaveCallbackInfo)

	//
	config.Port = firstNonZero(*flagPort, config.Port, 8080)

	//
	// database initialization

	etc.Database.CreateTableStruct(cTableSettings, itypes.RowSetting{})
	etc.Database.CreateTableStruct(cTableUsers, itypes.RowUser{})
	etc.Database.CreateTableStruct(cTableChannels, itypes.RowChannel{})
	etc.Database.CreateTableStruct(cTableRoles, itypes.RowRole{})
	etc.Database.CreateTableStruct(cTableChannelRolePerms, itypes.RowChannelRolePerms{})

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
	roleCache["o"] = itypes.RowRole{
		0, "o", 0, "Owner", "", pa, pa,
	}

	//
	// load roles into local cache

	for _, item := range queryAllRoles() {
		roleCache[item.UUID] = item
	}

	//
	// setup graceful stop

	util.RunOnClose(func() {
		util.Log("Gracefully shutting down...")

		util.Log("Saving database to disk")
		etc.Database.Close()

		util.Log("Closing all remaining active WebSocket connections")
		for _, item := range wsConnCache {
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
		writeAPIResponse(r, w, true, http.StatusOK, listToArray(connected))
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
		cv, ok := wsConnCache[user.UUID]
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
		conn, _ := wsUpgrader.Upgrade(w, r, nil)
		perms := calculateUserPermissions(user)
		wsConnCache[user.UUID] = itypes.ConnCacheValue{conn, user, perms}

		// connect
		if !listHas(connected, user.UUID) {
			connected.PushBack(user.UUID)
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
			smsg := string(msg)
			broadcastMessage(map[string]string{
				"type":    "message",
				"in":      smsg[:32],
				"from":    user.UUID,
				"message": smsg[32:],
			})
		}
		// disconnect
		if listHas(connected, user.UUID) {
			delete(wsConnCache, user.UUID)
			listRemove(connected, user.UUID)
			broadcastMessage(map[string]string{
				"type": "user-disconnect",
				"user": user.UUID,
			})
		}
	})

	//
	// start server

	if !util.IsPortAvailable(config.Port) {
		util.DieOnError(
			E(F("Binding to port %d failed.", config.Port)),
			"It may be taken or you may not have permission to. Aborting!",
		)
		return
	}

	p := strconv.Itoa(config.Port)
	util.Log("Initialization complete. Starting server on port " + p)
	http.ListenAndServe(":"+p, nil)
}
