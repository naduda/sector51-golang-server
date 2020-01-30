package apiserver

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/naduda/sector51-golang/internal/app/apiserver/httputils"
)

func (s *Server) handleUploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}

		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		if err = saveFile("./configs/"+handler.Filename, file); err != nil {
			httputils.SendError(w, http.StatusInternalServerError, err)
			return
		}

		httputils.Respond(w, http.StatusOK, nil)
	}
}

func saveFile(fn string, file multipart.File) error {
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return err
}
