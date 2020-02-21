package backup

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat = "20060102150405"
)

type Backup struct {
	logger              *logrus.Logger
	Period              int    `json:"period"`
	Count               int    `json:"count"`
	Folder              string `json:"dump_folder"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	UserName            string `json:"user"`
	Password            string `json:"password"`
	DbName              string `json:"dbName"`
	GoogleDriveFolderId string `json:"googleDriveFolderId"`
}

func New(logger *logrus.Logger) *Backup {
	period, err := strconv.Atoi(os.Getenv("BACKUP_PERIOD_SEC"))
	if err != nil {
		period = 15
	}
	count, err := strconv.Atoi(os.Getenv("BACKUP_COUNT"))
	if err != nil {
		count = 3
	}

	result := &Backup{
		logger: logger,
		Port:   5432,
		Folder: os.Getenv("DUMP_FOLDER"),
		Period: period,
		Count:  count,
	}

	if id, err := result.CreateDirIfNotExists(); err != nil {
		logger.Error(err.Error())
	} else {
		result.GoogleDriveFolderId = id
	}

	connStr := os.Getenv("CONNECTION_STR")
	connStr = strings.Split(connStr, "//")[1]
	arr := strings.Split(connStr, "@")

	tmp := strings.Split(arr[0], ":")
	result.UserName = tmp[0]
	result.Password = tmp[1]

	arr = strings.Split(arr[1], "/")
	tmp = strings.Split(arr[0], ":")
	result.Host = tmp[0]
	if len(tmp) > 1 {
		result.Port, _ = strconv.Atoi(tmp[1])
	}

	tmp = strings.Split(arr[1], "?")
	result.DbName = tmp[0]

	if _, err := os.Stat(result.Folder); os.IsNotExist(err) {
		if err = os.Mkdir(result.Folder, os.ModePerm); err != nil {
			logger.Error(err.Error())
		}
	}

	return result
}

func (b Backup) CreateDump() error {
	filename := fmt.Sprintf("%s/db.dump", b.Folder)

	args := []string{
		fmt.Sprintf("--port=%d", b.Port),
		fmt.Sprintf("--host=%s", b.Host),
		fmt.Sprintf("--username=%s", b.UserName),
		fmt.Sprintf("--dbname=%s", b.DbName),
		fmt.Sprintf("--file=%s", filename),
		"--format=custom",
		"--clean",
	}

	cmd := exec.Command("pg_dump", args...)
	if err := cmd.Start(); err != nil {
		b.logger.Error(err.Error())
		return err
	}

	return cmd.Wait()
}

// Restore ...
func (b Backup) Restore(filename string) error {
	args := []string{
		fmt.Sprintf("--port=%d", b.Port),
		fmt.Sprintf("--host=%s", b.Host),
		fmt.Sprintf("--username=%s", b.UserName),
		fmt.Sprintf("--dbname=%s", b.DbName),
		"--clean",
		filename,
	}

	cmd := exec.Command("pg_restore", args...)
	if err := cmd.Start(); err != nil {
		b.logger.Error(err.Error())
		return err
	}

	return cmd.Wait()
}

func (b *Backup) handleError(err error) {
	b.logger.Error(err.Error())
	time.Sleep(1 * time.Minute)
}

func (b *Backup) Start() {
	for {
		if b.GoogleDriveFolderId == "" {
			if id, err := b.CreateDirIfNotExists(); err != nil {
				b.handleError(err)
				continue
			} else {
				b.GoogleDriveFolderId = id
			}
		}

		if err := clearFolder(b.Folder); err != nil {
			b.handleError(err)
			continue
		}

		if err := b.CreateDumpAndUpload(); err != nil {
			b.handleError(err)
			continue
		}

		time.Sleep(time.Duration(b.Period) * time.Second)
	}
}

func (b Backup) CreateDumpAndUpload() (err error) {
	for i := 0; i < 5; i++ {
		if err = b.CreateDump(); err != nil {
			b.handleError(err)
			continue
		} else if err = b.Upload(); err != nil {
			b.handleError(err)
			continue
		}
		break
	}
	return
}

func clearFolder(folder string) error {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		if err := os.Mkdir(folder, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, f := range files {
		path := fmt.Sprintf("%s/%s", folder, f.Name())
		if err = os.Remove(path); err != nil {
			return err
		}
	}

	return nil
}

// ClearFolder ...
func ClearFolder(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// FindFileByExt ...
func FindFileByExt(dir, ext string) (string, error) {
	d, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return "", err
	}
	ext = strings.ToLower(ext)
	for _, name := range names {
		fmt.Println(name)
		if strings.HasSuffix(strings.ToLower(name), ext) {
			fmt.Printf("Name: %s\n\n\n%s", name, filepath.Join(dir, name))
			return filepath.Join(dir, name), nil
		}
	}
	return "", errors.New("file not found")
}
