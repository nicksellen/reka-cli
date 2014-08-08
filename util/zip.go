package util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func EmptyBuffer() *bytes.Buffer {
	return new(bytes.Buffer)
}

func Zip(root string) (*bytes.Buffer, error) {

	finfo, err := os.Stat(root)
	if err != nil {
		return EmptyBuffer(), fmt.Errorf("%s does not exist\n", root)
	}

	if finfo.IsDir() && !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		relPath := strings.TrimPrefix(strings.TrimPrefix(path, filepath.Dir(root)), "/")

		if strings.HasPrefix(info.Name(), ".") && len(info.Name()) > 1 {
			if info.IsDir() {
				//fmt.Printf("skipping %s\n", relPath)
				return filepath.SkipDir
			} else {
				//fmt.Printf("skipping %s\n", relPath)
				return nil
			}
		}

		if !info.IsDir() {
			//defer fmt.Printf("+ %s\n", relPath)

			header := &zip.FileHeader{
				Name:   relPath,
				Method: zip.Deflate,
			}

			header.SetModTime(info.ModTime())

			f, err := w.CreateHeader(header)

			if err != nil {
				log.Fatal(err)
			}

			file, err := os.Open(path)
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
