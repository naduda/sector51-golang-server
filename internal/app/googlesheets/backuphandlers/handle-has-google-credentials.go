package backuphandlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/googlesheets"
	"golang.org/x/oauth2"
)

// HandleHasGoogleCredentials ...
func HandleHasGoogleCredentials() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := make(map[string]interface{})
		config, err := googlesheets.GetGoogleConfig()
		result["hasFile"] = err == nil

		if err == nil {
			authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			result["authURL"] = authURL
			_, err = googlesheets.GetClient()
			result["exp"] = err != nil
			if err == nil {
				result["ok"] = true
			}
		}

		if gs, err := googlesheets.GetPageID(); err == nil {
			result["page"] = gs.PageID
		}

		httputils.Respond(w, http.StatusOK, result)
	}
}

// HandleCreateGoogleTokenFile ...
func HandleCreateGoogleTokenFile() http.HandlerFunc {
	type request struct {
		Code string `json:"code"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		config, err := googlesheets.GetGoogleConfig()
		if err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		tok, err := config.Exchange(context.TODO(), req.Code)
		if err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		if err = googlesheets.SaveToken(tok); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}
