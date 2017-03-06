package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MohamedBassem/godge"
)

func zipCurrentDir() ([]byte, error) {
	dir := "."
	zipfile := new(bytes.Buffer)

	archive := zip.NewWriter(zipfile)

	_, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read file stats: %v", err)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			header.Method = zip.Deflate
		}

		header.Name = path

		if info.IsDir() {
			header.Name += "/"
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
			return fmt.Errorf("failed to read file: %v", err)
		}

		if _, err := io.Copy(writer, file); err != nil {
			return fmt.Errorf("failed to copy file: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk dir: %v", err)
	}

	if err := archive.Close(); err != nil {
		return nil, fmt.Errorf("failed to close archive: %v", err)
	}
	return zipfile.Bytes(), nil
}

func checkResponseError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var e godge.ErrorResponse
		json.NewDecoder(resp.Body).Decode(&e)
		return errors.New(e.Error)
	}
	return nil
}
