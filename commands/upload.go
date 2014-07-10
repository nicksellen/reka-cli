package commands

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reka/config"
	"strings"
	"sync"
)

func Upload(Args []string) {

	if len(Args) != 1 {
		println("please provide path to upload")
		os.Exit(1)
	}

	config, err := config.Load()

	if err != nil {
		fmt.Printf("could not %s\n", err)
		return
	}

	root := Args[0]

	finfo, err := os.Stat(root)
	if err != nil {
		fmt.Printf("%s does not exist\n", root)
		return
	}
	if finfo.IsDir() && !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	fmt.Printf("uploading files from %s\n", root)

	var wg sync.WaitGroup

	walkFunc := func(path string, info os.FileInfo, err error) error {

		if strings.HasPrefix(info.Name(), ".") && len(info.Name()) > 1 {
			if info.IsDir() {
				fmt.Printf("skipping %s\n", path)
				return filepath.SkipDir
			} else {
				fmt.Printf("skipping %s\n", path)
				return nil
			}
		}

		if !info.IsDir() {

			wg.Add(1)

			go func(path string) {

				base := strings.TrimPrefix(path, filepath.Dir(root))

				if strings.HasPrefix(base, "/") {
					base = strings.TrimPrefix(base, "/")
				}

				defer wg.Done()
				defer fmt.Printf(" %s\n", base)

				url := fmt.Sprintf("%s/files/%s", config.URL, base)

				req, err := CreateUploadRequest(url, path)

				if err != nil {
					log.Fatal(err)
				}
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				} else {
					body := &bytes.Buffer{}
					_, err := body.ReadFrom(resp.Body)
					if err != nil {
						log.Fatal(err)
					}
					resp.Body.Close()
				}

			}(path)
		}

		return nil
	}

	filepath.Walk(root, walkFunc)

	wg.Wait()
}

func CreateUploadRequest(uri string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
