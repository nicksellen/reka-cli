package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nicksellen/reka/config"
	"github.com/nicksellen/reka/util"
	"github.com/wsxiaoys/terminal/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type DeployOk struct {
	Name     string    `json:"name"`
	Networks []Network `json:"network"`
}

type Network struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Url      string `json:"url"`
}

func Deploy(Args []string) {

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

	spinner1 := NewSpinner()
	color.Printf("packaging ")

	buf, err := util.Zip(config.WorkDir, config.Ignore)
	spinner1.Done()
	color.Printf("@{g}✓\n")
	if err != nil {
		log.Fatal(err)
	}

	spinner2 := NewSpinner()

	color.Printf("deploying to @{!}%s@{|} ", deployment.Name)

	client := &http.Client{}
	req, err := http.NewRequest("POST", deployment.Url, buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/zip")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)

	/*
		resp, err := http.Post(deployment.Url, "application/zip", buf)
	*/
	spinner2.Done()

	if err != nil {
		color.Printf("@{r}✕ failed\n  @{.}%s@{|}\n", err)
		os.Exit(1)
	}

	code := resp.StatusCode

	var data map[string]interface{}

	switch {
	case code >= 200 && code < 300:
		var ok DeployOk
		json.Unmarshal(ReadBody(resp), &ok)

		var buffer bytes.Buffer
		color.Fprintf(&buffer, "@{g}✓ success@{|}")
		color.Fprintf(&buffer, "\n\n")

		color.Fprintf(&buffer, "          @{.}name@{|} @{!}%s@{|}\n", ok.Name)
		color.Fprintf(&buffer, "    @{.}deployment@{|} @{!}%s@{|}\n", deployment.Url)

		if len(ok.Networks) > 0 {
			color.Fprintf(&buffer, "\n  @{.}listening on@{|} ")
			for _, network := range ok.Networks {
				if network.Url != "" {
					color.Fprintf(&buffer, "@{!}%s@{|}\n", network.Url)
				} else {
					color.Fprintf(&buffer, "%s on port %d\n", network.Protocol, network.Port)
				}
				color.Fprintf(&buffer, "               ")
			}
		}
		color.Fprintf(&buffer, "\n")

		fmt.Print(buffer.String())
	case code >= 400 && code < 500:
		json.Unmarshal(ReadBody(resp), &data)
		color.Printf("@{r}✕ client error@{|}\n@{.}%s@{|}\n", data["message"])
	case code >= 500 && code < 600:
		json.Unmarshal(ReadBody(resp), &data)
		color.Printf("@{r}✕ error@{|}\n@{.}%s@{|}\n", data["message"])
	default:
		color.Printf("@{r}✕ %d error@{|}\n", code)
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

type Spinner struct {
	IsDone bool
	wg     sync.WaitGroup
}

func (spinner Spinner) Done() {
	spinner.IsDone = true
	spinner.wg.Wait()
}

func NewSpinner() Spinner {

	var wg sync.WaitGroup

	spinner := Spinner{
		wg:     wg,
		IsDone: false,
	}

	wg.Add(1)

	go func() {
		chars := [4]string{"\\", "|", "/", "-"}
		next := 0
		defer wg.Done()
		for !spinner.IsDone {
			os.Stdout.Write([]byte(chars[next]))
			os.Stdout.Write([]byte("\b"))
			time.Sleep(100 * time.Millisecond)
			next = next + 1
			if next > len(chars)-1 {
				next = 0
			}
		}
	}()

	return spinner
}
