package clients

import (
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/store"
	"net/http"
)

func HandleServices(repo store.ServiceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := repo.List()
		if err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}
		httputils.Respond(w, http.StatusOK, res)
	}
}
