package backup

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

// CreateDir inside google drive with if not Exists.
func (b Backup) CreateDirIfNotExists() (string, error) {
	srv, err := GetDriveService()
	if err != nil {
		return "", err
	}

	if r, err := srv.Files.List().Q("'root' in parents").Fields("files(id, name)").Do(); err != nil {
		return "", err
	} else {
		for _, i := range r.Files {
			if i.Name == b.Folder {
				return i.Id, nil
			}
		}
	}

	return b.CreateDir(srv, b.Folder)
}

// CreateDir inside google drive.
func (b Backup) CreateDir(srv *drive.Service, folderName string) (string, error) {
	d := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{},
	}
	dir, err := srv.Files.Create(d).Do()
	if err != nil {
		return "", err
	}
	return dir.Id, nil
}

// CreateFile into google drive.
func CreateFile(srv *drive.Service, name string, fileToUpload *os.File, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType:                     "application/tar+gzip",
		Name:                         name,
		Parents:                      []string{parentId},
		CopyRequiresWriterPermission: true,
		WritersCanShare:              true,
	}
	return srv.Files.Create(f).Media(fileToUpload).Do()
}

// Upload backup file
func (b *Backup) Upload() error {
	srv, err := GetDriveService()
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/db.dump", b.Folder)

	zipFile := fmt.Sprintf("%s.zip", filename)
	if err := ZipFiles(zipFile, []string{filename}); err != nil {
		return err
	}

	f, err := os.Open(zipFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := deleteBackups(srv, b.GoogleDriveFolderId, b.Count); err != nil {
		b.logger.Error(err.Error())
	}

	currentDate := time.Now().Format(DateFormat)
	filename = fmt.Sprintf("db_%s.zip", currentDate)

	_, err = CreateFile(srv, filename, f, b.GoogleDriveFolderId)
	return err
}

func (b *Backup) GetDumpList() ([]string, error) {
	srv, err := GetDriveService()
	if err != nil {
		return nil, err
	}

	folder := fmt.Sprintf("'%s' in parents", b.GoogleDriveFolderId)
	if r, err := srv.Files.List().Q(folder).Fields("files(id, name)").Do(); err == nil {
		files := r.Files
		var res []string
		for _, f := range files {
			b.logger.Info(f.Id, f.Name)
			res = append(res, f.Id)
		}
		return res, nil
	} else {
		return nil, err
	}
}

// Download backup file
func (b *Backup) Download(id string, filename string) error {
	srv, err := GetDriveService()
	if err != nil {
		return err
	}

	res, err := srv.Files.Get(id).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, body, 0644)
}

func deleteBackups(srv *drive.Service, folderId string, limit int) error {
	folder := fmt.Sprintf("'%s' in parents", folderId)
	if r, err := srv.Files.List().Q(folder).Fields("files(id, createdTime)").Do(); err == nil {
		files := r.Files
		countFiles := len(files)
		if countFiles < limit {
			return nil
		}
		sort.Slice(files[:], func(i, j int) bool {
			return files[i].CreatedTime < files[j].CreatedTime
		})

		for i := 0; i <= countFiles-limit; i++ {
			file := files[i]
			if err := srv.Files.Delete(file.Id).Do(); err != nil {
				return err
			}
		}
	}
	return nil
}
