package clients

import (
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/store"
	"net/http"
)

func HandleUserServices(repo store.ServiceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.URL.Query()["id"]
		if !ok {
			httputils.SendError(w, http.StatusBadRequest, store.ErrIllegalArgs)
		}
		res, err := repo.GetUserServices(userId[0])
		if err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}
		httputils.Respond(w, http.StatusOK, res)
	}
}
