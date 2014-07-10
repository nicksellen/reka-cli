package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reka/config"
)

func Undeploy(Args []string) {

	if len(Args) != 1 {
		println("please provide path to deploy from")
		os.Exit(1)
	}

	config, err := config.Load()

	client := &http.Client{}

	uuid := Args[0]

	url := fmt.Sprintf("%s/apps/%s", config.URL, uuid)

	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Do(req)

	fmt.Printf("undeployed %s!\n", uuid)

}
