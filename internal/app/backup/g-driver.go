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
	filename := fmt.Sprintf("%s/dump.sql", b.Folder)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := deleteBackups(srv, b.GoogleDriveFolderId, b.Count); err != nil {
		b.logger.Error(err.Error())
	}

	currentDate := time.Now().Format(DateFormat)
	filename = fmt.Sprintf("dump_%s.sql", currentDate)

	_, err = CreateFile(srv, filename, f, b.GoogleDriveFolderId)
	return err
}

// Download backup file
func (b *Backup) Download(id string) error {
	srv, err := GetDriveService()
	if err != nil {
		return err
	}

	res, err := srv.Files.Get(id).Download()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}

	fmt.Println(string(body))

	return err
}

func deleteBackups(srv *drive.Service, folderId string, limit int) error {
	folder := fmt.Sprintf("'%s' in parents", folderId)
	if r, err := srv.Files.List().Q(folder).Fields("files(id, name)").Do(); err == nil {
		files := r.Files
		countFiles := len(files)
		if countFiles < limit {
			return nil
		}
		sort.Slice(files[:], func(i, j int) bool {
			return files[i].Name < files[j].Name
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
