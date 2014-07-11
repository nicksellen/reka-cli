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

func Deploy(Args []string) {

	if len(Args) != 1 {
		println("please provide path to deploy")
		os.Exit(1)
	}

	spec := Args[0]

	config, err := config.Load()

	resp, err := http.PostForm(fmt.Sprintf("%s/apps", config.URL), url.Values{"spec": {spec}})

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

	var data map[string]interface{}

	err = json.Unmarshal(body.Bytes(), &data)

	fmt.Println(data["message"])

}
