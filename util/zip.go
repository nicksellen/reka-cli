package util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func EmptyBuffer() *bytes.Buffer {
	return new(bytes.Buffer)
}

func Zip(root string, ignore []string) (*bytes.Buffer, error) {

	finfo, err := os.Stat(root)
	if err != nil {
		return EmptyBuffer(), fmt.Errorf("%s does not exist\n", root)
	}

	if finfo.IsDir() && !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	filepath.Walk(root, func(p string, info os.FileInfo, err error) error {

		relPath := strings.TrimPrefix(strings.TrimPrefix(p, filepath.Dir(root)), "/")

		if info == nil {
			return nil
		}

		if strings.HasPrefix(info.Name(), ".") && len(info.Name()) > 1 {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		for _, ignorePattern := range ignore {
			matched, _ := path.Match(ignorePattern, relPath)
			if matched {
				return nil
			}
		}

		if !info.IsDir() {

			header := &zip.FileHeader{
				Name:   relPath,
				Method: zip.Deflate,
			}

			header.SetModTime(info.ModTime())

			f, err := w.CreateHeader(header)

			if err != nil {
				log.Fatal(err)
			}

			file, err := os.Open(p)
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(f, bufio.NewReader(file))
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	err = w.Close()
	if err != nil {
		return EmptyBuffer(), err
	}

	return buf, nil

}
