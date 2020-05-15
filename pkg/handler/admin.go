package handler

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/ws"
)

// AuditsCsv handles /api/admin/audits.csv
func AuditsCsv(w http.ResponseWriter, r *http.Request) {
	_, user, err := apiBootstrapRequireLogin(r, w, http.MethodGet, true)
	if err != nil {
		return
	}
	usp := ws.UserPerms{}.From(user)
	if !usp.ViewAudits {
		writeAPIResponse(r, w, false, http.StatusForbidden, "action requires the view_audits permission")
		return
	}
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
