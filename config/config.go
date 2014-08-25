package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RekaConfig struct {
	WorkDir       string            `json:"work"`
	DefaultServer string            `json:"default-server"`
	Servers       map[string]Server `json:"servers"`
}

type Server struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (config RekaConfig) GetServer(name string) (Server, error) {
	server, ok := config.Servers[name]
	if !ok {
		var keys []string
		for k := range config.Servers {
			keys = append(keys, k)
		}
		return *new(Server), fmt.Errorf("no server called %s (valid servers: %s)\n", name, keys)
	}
	return server, nil
}

func Load() (RekaConfig, error) {
	return Walk()
}

func TrimmedContents(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(bytes))
}

func Walk() (RekaConfig, error) {
	root := ".reka/config/"
	var config RekaConfig

	config.WorkDir = "."

	config.Servers = make(map[string]Server)

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := strings.TrimPrefix(path, root)

			if strings.HasPrefix(key, "servers/") {
				name := strings.TrimPrefix(key, "servers/")
				config.Servers[name] = Server{
					Name: name,
					URL:  TrimmedContents(path),
				}
			} else if key == "default-server" {
				config.DefaultServer = TrimmedContents(path)
			} else {
				fmt.Printf("walked unused '%s'\n", key)
			}

		}
		return nil
	})
	return config, nil
}
