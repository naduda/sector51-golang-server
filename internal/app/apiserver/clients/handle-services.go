package clients

import (
	"encoding/json"
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
	"net/http"
)

func HandleServices(repo store.ServiceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			res, err := repo.List()
			if err != nil {
				httputils.SendError(w, http.StatusInternalServerError, err)
				return
			}
			httputils.Respond(w, http.StatusOK, res)
		} else if r.Method == "PUT" {
			req := &model.Service{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				httputils.SendError(w, http.StatusBadRequest, err)
				return
			}

			if err := repo.UpdateService(*req); err != nil {
				httputils.SendError(w, http.StatusBadRequest, err)
				return
			}
			httputils.Respond(w, http.StatusOK, nil)
		}
	}
}
