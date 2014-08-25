package commands

import (
	"encoding/json"
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"io/ioutil"
	"log"
	"net/http"
	"reka/config"
	"reka/util"
)

func Push(Args []string) {

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	name := ""

	if len(Args) == 0 {
		name = config.DefaultServer
		if name == "" {
			log.Fatal("please provide reka server name, or setup a default server: echo name > .reka/config/default-server")
		}
	} else {
		name = Args[0]
	}

	server, err := config.GetServer(name)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := util.Zip(config.WorkDir)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/apps/%s", server.URL, config.Identity)

	color.Printf("pushing @{!}%s@{|} to %s... ", config.Identity, server.Name)

	resp, err := http.Post(url, "application/zip", buf)
	if err != nil {
		color.Printf("@{r}✕ failed\n")
		log.Fatal(err)
	}
	code := resp.StatusCode
	ct := resp.Header.Get("Content-Type")
	if ct == "application/json" {
		var data map[string]interface{}
		json.Unmarshal(ReadBody(resp), &data)
		color.Printf("@{r}✕ server error@{|}\n%s\n", data["message"])
	} else {
		switch {
		case code >= 200 && code < 300:
			color.Printf("@{g}✓ %s\n", ReadBody(resp))
		case code >= 400 && code < 500:
			color.Printf("@{r}✕ client error@{|}\n%s\n", string(ReadBody(resp)))
		case code >= 500 && code < 600:
			color.Printf("@{r}✕ server error@{|}\n%s\n", string(ReadBody(resp)))
		default:
			color.Printf("@{r}✕ %d error@{|}\n", code)
		}
	}

}

func ReadBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
