package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go.etc/htp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Gauge setup
var (
	GaugeUsrTotal     = newGauge("users_total")
	GaugeUsrBy        = newGaugeLabeled("users_by", "stat")
	GaugeChanTotal    = newGauge("channels_total")
	GaugeChanMsgTotal = newGaugeLabeled("messages_total", "channel")
	GaugeChanMsgEdits = newGaugeLabeled("messages_edits", "channel")
	GaugeChanMsgDels  = newGaugeLabeled("messages_deletes", "channel")
	GaugeRoleTotal    = newGauge("roles_total")
	GaugeInvTotal     = newGauge("invites_total")
	GaugeInvUses      = newGaugeLabeled("invites_uses", "code")
	GaugeAudTotal     = newGauge("audits_total")
	GaugeAudBy        = newGaugeLabeled("audits_by", "action")
)

func refresh() {
	GaugeUsrTotal.Set(getPropInt("count_users"))
	GaugeChanTotal.Set(getPropInt("count_channels"))
	GaugeRoleTotal.Set(getPropInt("count_roles"))
	GaugeInvTotal.Set(getPropInt("count_invites"))
	GaugeAudTotal.Set(getPropInt("count_audits"))

	GaugeUsrBy.With(prometheus.Labels{"stat": "member"}).Set(getPropInt("count_users_members"))
	GaugeUsrBy.With(prometheus.Labels{"stat": "online"}).Set(float64(ws.OnlineUserCount()))
	GaugeUsrBy.With(prometheus.Labels{"stat": "banned"}).Set(getPropInt("count_users_banned"))

	for _, item := range (db.Channel{}.All()) {
		uid := item.UUID.String()
		GaugeChanMsgTotal.With(prometheus.Labels{"channel": uid}).Set(getPropInt("count_messages_" + uid))
		GaugeChanMsgEdits.With(prometheus.Labels{"channel": uid}).Set(getPropInt("count_messages_" + uid + "_edits"))
		GaugeChanMsgDels.With(prometheus.Labels{"channel": uid}).Set(getPropInt("count_messages_" + uid + "_deletes"))
	}
	for _, item := range (db.Invite{}.All()) {
		GaugeInvUses.With(prometheus.Labels{"code": item.Code}).Set(float64(item.Uses))
	}
	for i := 1; i < db.ActionLen(); i++ {
		is := strconv.Itoa(i)
		GaugeAudBy.With(prometheus.Labels{"action": is}).Set(getPropInt("count_audits_action_" + is))
	}
}

// Init sets initial values and creates needed variables
func Init() {
	go func() {
		for {
			refresh()
			time.Sleep(time.Second * 5)
		}
	}()
}

// Handler returns the prometheus http.HandlerFunc
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := htp.GetController(r)
		withKey := r.Header.Get("Authorization") == "Bearer "+db.Props.Get("prometheus_key")
		if !withKey {
			u := controls.GetMemberUser(c, r, w)
			c.Assert(u.HasRole("o"), "403: resource requires Authorization or to be server owner to access")
		}
		promhttp.Handler().ServeHTTP(w, r)
	}
}
