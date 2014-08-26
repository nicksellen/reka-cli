package commands

import (
	"encoding/json"
	"github.com/wsxiaoys/terminal/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reka/config"
	"reka/util"
	"sync"
	"time"
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

	spinner1 := NewSpinner()
	color.Printf("packing ")

	buf, err := util.Zip(config.WorkDir)
	spinner1.Done()
	color.Printf("@{g}✓\n")
	if err != nil {
		log.Fatal(err)
	}

	spinner2 := NewSpinner()

	color.Printf("pushing to @{!}%s@{|} ", server.Name)

	resp, err := http.Post(server.URL, "application/zip", buf)
	spinner2.Done()

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
