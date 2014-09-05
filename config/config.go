package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RekaConfig struct {
	WorkDir           string                `json:"work"`
	DefaultDeployment string                `json:"default-deployment"`
	Deployments       map[string]Deployment `json:"deployments"`
	Ignore            []string              `json:"ignore"`
}

type Deployment struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (config RekaConfig) GetDeployment(name string) (Deployment, error) {
	deployment, ok := config.Deployments[name]
	if !ok {
		var keys []string
		for k := range config.Deployments {
			keys = append(keys, k)
		}
		return *new(Deployment), fmt.Errorf("no deployment %s, try one of %s\n", name, keys)
	}
	return deployment, nil
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
	config.Deployments = make(map[string]Deployment)
	config.Ignore = []string{}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := strings.TrimPrefix(path, root)
			if strings.HasPrefix(key, "deployments/") {
				name := strings.TrimPrefix(key, "deployments/")
				config.Deployments[name] = Deployment{
					Name: name,
					Url:  TrimmedContents(path),
				}
			} else if key == "default-deployment" {
				config.DefaultDeployment = TrimmedContents(path)
			} else if strings.HasPrefix(key, "servers/") {
				fmt.Printf("please move 'server' to 'deployment' config at %s\n", path)
				name := strings.TrimPrefix(key, "servers/")
				config.Deployments[name] = Deployment{
					Name: name,
					Url:  TrimmedContents(path),
				}
			} else if key == "default-server" {
				fmt.Printf("please move 'default-server' to 'default-deployment' config at %s\n", path)
				config.DefaultDeployment = TrimmedContents(path)
			} else if key == "ignore" {
				config.Ignore = readLines(path)
			} else {
				fmt.Printf("deprecated config at %s\n", path)
			}

		}
		return nil
	})
	return config, nil
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return lines
}
