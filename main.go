package main

import (
	"log"

	oauth2 "github.com/nektro/go.oauth2"

	. "github.com/nektro/go-util/alias"

	_ "github.com/nektro/mantle/statik"
)

var (
	config      *Config
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
}
