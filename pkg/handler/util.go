package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"
	sdrie "github.com/nektro/go.sdrie"
)

var (
	badgeCache = sdrie.New()
)

// Init sets up this package
func Init() {
	etc.HtpErrCb = func(r *http.Request, w http.ResponseWriter, good bool, code int, data string) {
		writeAPIResponse(r, w, good, code, data)
	}
}

func writeAPIResponse(r *http.Request, w http.ResponseWriter, good bool, status int, message interface{}) {
	resp := map[string]interface{}{
		"success": good,
		"message": message,
	}
	w.Header().Add("content-type", "application/json")
	dat, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(dat))
}

func hGrabInt(s string) (string, int64, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	return s, n, err
}

func hBadge(w http.ResponseWriter, r *http.Request, l, m, c string) {
	l = strings.ReplaceAll(l, " ", "_")
	m = strings.ReplaceAll(m, " ", "_")
	k := l + "-" + m + "-" + c
	w.Header().Add("content-type", "image/svg+xml")
	if badgeCache.Has(k) {
		fmt.Fprintln(w, string(badgeCache.Get(k).([]byte)))
		return
	}
	req, _ := http.NewRequest(http.MethodGet, "https://img.shields.io/badge/"+k, nil)
	bys, err := util.DoHttpFetch(req)
	if err != nil {
		return
	}
	badgeCache.Set(k, bys, time.Minute*10)
	fmt.Fprintln(w, string(bys))
}

func hBetween(x, l, h int) bool {
	if x < l {
		return false
	}
	if x > h {
		return false
	}
	return true
}
