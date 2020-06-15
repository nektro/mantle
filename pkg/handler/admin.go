package handler

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/handler/controls"
	"github.com/nektro/mantle/pkg/ws"

	"github.com/nektro/go.etc/htp"
)

// AuditsCsv handles /api/admin/audits.csv
func AuditsCsv(w http.ResponseWriter, r *http.Request) {
	c := htp.GetController(r)
	user := controls.GetMemberUser(c, r, w)
	usp := ws.UserPerms{}.From(user)
	c.Assert(usp.ViewAudits, "403: action requires the view_audits permission")

	w.Header().Add("content-type", "text/csv")
	cw := csv.NewWriter(w)
	for _, item := range (db.Audit{}.All()) {
		cw.Write([]string{
			strconv.FormatInt(item.ID, 10),
			item.UUID,
			item.CreatedOn.String(),
			item.Action.String(),
			item.Agent,
			item.Affected,
			item.Key,
			item.Value,
		})
		cw.Flush()
	}
}
