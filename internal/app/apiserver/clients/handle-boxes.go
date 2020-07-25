package clients

import (
	"encoding/json"
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
	"net/http"
)

func HandleBoxes(repo store.BoxRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := repo.List()
		if err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}
		httputils.Respond(w, http.StatusOK, res)
	}
}

func HandleBox(repo store.BoxRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &model.Box{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		if err := repo.Update(*req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}
