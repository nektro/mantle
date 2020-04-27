package itypes

import (
	oauth2 "github.com/nektro/go.oauth2"
)

type Config struct {
	Port      int
	Clients   []oauth2.AppConf  `json:"clients"`
	Providers []oauth2.Provider `json:"providers"`
	Themes    []string          `json:"themes"`
}
