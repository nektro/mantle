package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	oauth2 "github.com/nektro/go.oauth2"

	. "github.com/nektro/go-util/alias"

	_ "github.com/nektro/mantle/statik"
)

var (
	config      *Config
	roleCache   = map[string]RowRole{}
)

type Config struct {
	Version   int               `json:"version"`
	Port      int               `json:"port"`
	Clients   []oauth2.AppConf  `json:"clients"`
	Providers []oauth2.Provider `json:"providers"`
}

func main() {
	log.Println("Welcome to " + Name + ".")

	//
	etc.Init("mantle", &config, "./invite", helperSaveCallbackInfo)

	//
	// database initialization

	etc.Database.CreateTableStruct(cTableSettings, RowSetting{})
	etc.Database.CreateTableStruct(cTableUsers, RowUser{})
	etc.Database.CreateTableStruct(cTableChannels, RowChannel{})
	etc.Database.CreateTableStruct(cTableRoles, RowRole{})
	etc.Database.CreateTableStruct(cTableChannelRolePerms, RowChannelRolePerms{})

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

	pa := PermAllow
	roleCache["o"] = RowRole{
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
		log.Println("Gracefully shutting down...")

		log.Println("Saving database to disk")
		etc.Database.Close()

		log.Println("Done")
		os.Exit(0)
	})

	//
	// create http service

	http.HandleFunc("/invite", func(w http.ResponseWriter, r *http.Request) {
		_, user, _ := apiBootstrapRequireLogin(r, w, http.MethodGet, false)

		if props.Get("public") == "true" {
			if user.IsMember == false {
				etc.Database.Build().Up(cTableUsers, "is_member", "1").Wh("uuid", user.UUID).Exe()
				log.Println("[user-join]", F("User %s just became a member and joined the server", user.UUID))
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

	//
	// start server

	p := strconv.Itoa(config.Port)
	log.Println("Initialization complete. Starting server on port " + p)
	http.ListenAndServe(":"+p, nil)
}
