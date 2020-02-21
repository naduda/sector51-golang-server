package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzipping
func Unzipping(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		//noinspection ALL
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		//noinspection ALL
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func Zipping(src, dest string) error {
	zipfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	//noinspection ALL
	defer archive.Close()

	info, err := os.Stat(src)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, src))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}
