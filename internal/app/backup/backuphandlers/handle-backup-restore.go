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

		err = backup.ClearFolder("/tmp")
		if err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		zipFile := "/tmp/db.zip"
		if err = b.Download(list[0], zipFile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		dumpFile := "/tmp"
		if err := backup.Unzipping(zipFile, dumpFile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		dumpFile, err = backup.FindFileByExt("/tmp", ".dump")
		if err := backup.Unzipping(zipFile, dumpFile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		if err = b.Restore("/tmp/" + dumpFile); err != nil {
			httputils.SendError(w, http.StatusBadRequest, err)
			return
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}
