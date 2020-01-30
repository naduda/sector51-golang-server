package backuphandlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/googlesheets"
)

// HandleSaveSheetID ...
func HandleSaveSheetID() http.HandlerFunc {
	type request struct {
		PageID string `json:"pageId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		data, _ := json.Marshal(req)
		if err := ioutil.WriteFile(googlesheets.PageFile, data, 0644); err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}
