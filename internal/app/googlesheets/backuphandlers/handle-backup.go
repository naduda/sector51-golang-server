package backuphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/googlesheets"
)

// HandleBackup ...
func HandleBackup() http.HandlerFunc {
	type request struct {
		SheetID string `json:"sheetId"`
	}
	e := googlesheets.NewSheetExporter()

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		if err := e.Backup(req.SheetID); err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}
