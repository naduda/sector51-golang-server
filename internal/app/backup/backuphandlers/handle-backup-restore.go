package backuphandlers

import (
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

		logger.Error("Downloading...", list[0])
		dumpfile := "/tmp/dump.sql"
		if err = b.Download(list[0], dumpfile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		if err = b.Restore(dumpfile); err != nil {
			logger.Error("Err 3", err.Error())
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}
