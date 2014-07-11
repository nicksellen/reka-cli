package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reka/config"
)

func Redeploy(Args []string) {

	if len(Args) != 1 {
		println("please provide the id to redeploy")
		os.Exit(1)
	}

	config, err := config.Load()

	client := &http.Client{}

	identity := Args[0]

	url := fmt.Sprintf("%s/apps/%s/redeploy", config.URL, identity)

	req, err := http.NewRequest("PUT", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Do(req)

	fmt.Printf("redeployed %s!\n", identity)

}
