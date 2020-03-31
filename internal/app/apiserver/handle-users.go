package apiserver

import (
	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
	"github.com/naduda/sector51-golang/internal/app/model"
	"net/http"
	"regexp"
	"strings"
)

func (s *Server) handleUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, ok := r.URL.Query()["f"]
		f := ""
		if ok {
			f = strings.ToLower(filter[0])
		}

		res, err := s.store.User().FindAll()
		if err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}

		result := []model.User{}
		if f == "" {
			for _, v := range res {
				if v.Active {
					result = append(result, v)
				} else {
					break
				}
			}
			httputils.Respond(w, http.StatusOK, result)
			return
		}

		f = strings.ToLower(f)
		for _, v := range res {
			var re = regexp.MustCompile(`\D`)
			criteria := strings.ToLower(re.ReplaceAllString(v.Phone, ``))
			if strings.Contains(criteria, f) {
				result = append(result, v)
				continue
			}

			criteria = strings.ToLower(v.Name)
			if strings.Contains(criteria, f) {
				result = append(result, v)
				continue
			}

			criteria = strings.ToLower(v.Surname)
			if strings.Contains(criteria, f) {
				result = append(result, v)
				continue
			}

			if strings.Contains(v.Card, f) {
				result = append(result, v)
				continue
			}
		}
		httputils.Respond(w, http.StatusOK, result)
	}
}
