package godge

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
)

func randomString(n int) string {
	letterBytes := "abcdefghijklmnopkrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func unzipToTmpDir(b []byte) (string, error) {
	tdir, err := ioutil.TempDir("", "godge")
	if err != nil {
		return "", fmt.Errorf("failed to create a tmp dir: %v", err)
	}
	r := bytes.NewReader(b)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return "", fmt.Errorf("failed to create a new zip reader: %v", err)
	}

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(tdir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range zr.File {
		if err := extractAndWriteFile(f); err != nil {
			return "", err
		}
	}
	return tdir, nil
}
