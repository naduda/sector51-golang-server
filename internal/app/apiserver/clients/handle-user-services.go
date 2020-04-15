package clients

import (
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
	"github.com/naduda/sector51-golang/internal/app/store"
	"net/http"
	"time"
)

func HandleUserServices(repo store.ServiceRepository) http.HandlerFunc {
	type UserServiceDTO struct {
		IdService int    `json:"idService"`
		IdUser    string `json:"idUser"`
		DtBeg     int64  `json:"dtBeg"`
		DtEnd     int64  `json:"dtEnd"`
		Value     string `json:"value"`
	}

	toDTO := func(services []model.UserService) []UserServiceDTO {
		res := make([]UserServiceDTO, len(services))

		for i, s := range services {
			res[i] = UserServiceDTO{
				IdService: s.IdService,
				IdUser:    s.IdUser,
				DtBeg:     s.DtBeg.UnixNano() / int64(time.Millisecond),
				DtEnd:     s.DtEnd.UnixNano() / int64(time.Millisecond),
				Value:     s.Value,
			}
		}

		return res
	}

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

		httputils.Respond(w, http.StatusOK, toDTO(res))
	}
}
