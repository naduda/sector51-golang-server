package backup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat = "20060102150405"
	FileJson   = "configs/backup.json"
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

func NewBackup(connStr string, logger *logrus.Logger) *Backup {
	if result, err := backupFromFile(); err == nil {
		result.logger = logger
		return result
	}

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

	if err := result.saveToJson(); err != nil {
		logger.Error("Can't save backup json")
	}
	return result
}

func backupFromFile() (*Backup, error) {
	f, err := os.Open(FileJson)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &Backup{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (b *Backup) saveToJson() error {
	f, err := os.OpenFile(FileJson, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b)
}

func (b Backup) dumpCommand() *exec.Cmd {
	args := []string{
		fmt.Sprintf("--port=%d", b.Port),
		fmt.Sprintf("--host=%s", b.Host),
		fmt.Sprintf("--username=%s", b.UserName),
		fmt.Sprintf("--dbname=%s", b.DbName),
	}

	cmd := exec.Command("pg_dump", args...)

	//cmdTemplate := "pg_dump --port=%d --host=%s --username=%s, --dbname=%s > %s/dump.sql"
	//command := fmt.Sprintf(cmdTemplate, b.Port, b.Host, b.UserName, b.DbName, b.Folder)
	//b.logger.Info(command)
	return cmd
}

func (b Backup) CreateDump() error {
	cmd := b.dumpCommand()

	stdout, err := cmd.StdoutPipe()
	var out bytes.Buffer
	cmd.Stderr = &out
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		b.logger.Error(err.Error())
		return err
	}

	b.logger.Debug("killing...")
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	//case <-time.After(1 * time.Millisecond):
	//	if err := cmd.Process.Kill(); err != nil {
	//		b.logger.Error("failed to kill process: ", err)
	//	}
	//	b.logger.Info("process killed as timeout reached")
	case err := <-done:
		if err != nil {
			b.logger.Error("process finished with error = %v", err)
		}
		b.logger.Info("process finished successfully")
		if err := cmd.Process.Kill(); err != nil {
			b.logger.Error("failed to kill process: ", err)
		}
	}

	bytesArray, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/dump.sql", b.Folder)
	return ioutil.WriteFile(filename, bytesArray, os.ModePerm)
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
				if err = b.saveToJson(); err != nil {
					b.handleError(err)
					continue
				}
			}
		}

		if err := clearFolder(b.Folder); err != nil {
			b.handleError(err)
			continue
		}

		if err := b.CreateDump(); err != nil {
			b.handleError(err)
			continue
		} else if err = b.Upload(); err != nil {
			b.handleError(err)
			continue
		}

		time.Sleep(time.Duration(b.Period) * time.Second)
	}
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
