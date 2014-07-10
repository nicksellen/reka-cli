package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type CLIConfig struct {
	URL string `json:"url"`
}

func Load() (CLIConfig, error) {
	var filename = ".reka.json"
	var config CLIConfig
	if _, err := os.Stat(filename); err == nil {
		var data, err = ioutil.ReadFile(filename)
		if err != nil {
			return config, err
		}
		json.Unmarshal(data, &config)
	} else {
		println("missing config file", filename)
		os.Exit(1)
	}
	return config, nil
}
