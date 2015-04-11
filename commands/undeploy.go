package commands

import (
	"github.com/nicksellen/reka/config"
	"github.com/wsxiaoys/terminal/color"
	"log"
	"net/http"
	"os"
)

func Undeploy(Args []string) {

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	name := ""

	if len(Args) == 0 {
		name = config.DefaultDeployment
		if name == "" {
			log.Fatal("please provide reka deployment name, or setup a default: echo name > .reka/config/default-deployment")
		}
	} else {
		name = Args[0]
	}

	deployment, err := config.GetDeployment(name)
	if err != nil {
		log.Fatal(err)
	}

	color.Printf("undeploying from @{!}%s@{|} ", deployment.Name)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", deployment.Url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		color.Printf("@{r}✕ failed\n  @{.}%s@{|}\n", err)
		os.Exit(1)
	}

	code := resp.StatusCode

	switch {
	case code >= 200 && code < 300:
		color.Print("@{g}✓@{|}\n")
	default:
		color.Printf("@{r}✕ %d error@{|}\n", code)
	}

}
