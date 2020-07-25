package backuphandlers

import (
	"encoding/json"
	"github.com/naduda/sector51-golang/internal/app/backup"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
)

// HandleBackup ...
func HandleBackup(logger *logrus.Logger) http.HandlerFunc {
	b := backup.New(logger)

	return func(w http.ResponseWriter, r *http.Request) {
		if err := b.CreateDumpAndUpload(); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}
		httputils.Respond(w, http.StatusOK, nil)
	}
}

// HandleRestore ...
func HandleRestore(logger *logrus.Logger) http.HandlerFunc {
	b := backup.New(logger)

	return func(w http.ResponseWriter, r *http.Request) {
		list, err := b.GetDumpList()
		if err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		err = backup.ClearFolder("/tmp")
		if err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		dumpFile := "/tmp/db.dump"
		if err = b.Download(list[0], dumpFile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		if err = b.Restore(dumpFile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}

// HandleRestore2 ...
func HandleRestore2(logger *logrus.Logger) http.HandlerFunc {
	b := backup.New(logger)

	type NamedDTO struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &NamedDTO{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		dumpFile := "./configs/" + req.Name
		if err := b.Restore(dumpFile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}
