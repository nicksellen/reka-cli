package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"reka/config"
)

func Validate(Args []string) {

	if len(Args) != 1 {
		println("please provide path to validate from")
		os.Exit(1)
	}

	spec := Args[0]

	config, err := config.Load()

	resp, err := http.PostForm(fmt.Sprintf("%s/validate", config.URL), url.Values{"spec": {spec}})

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body := &bytes.Buffer{}

	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//var data map[string]interface{}

	type ValidationError struct {
		Message     string `json:"message"`
		Source      string `json:"source"`
		Linenumbers string `json:"linenumbers"`
	}

	type ValidationResult struct {
		Message string            `json:"message"`
		Errors  []ValidationError `json:"errors"`
	}

	var data ValidationResult

	err = json.Unmarshal(body.Bytes(), &data)

	if data.Errors != nil {
		fmt.Printf("%d error(s)\n", len(data.Errors))
		for _, error := range data.Errors {
			fmt.Println(error.Linenumbers, error.Message)
		}
	}

	fmt.Println(data.Message)

}
