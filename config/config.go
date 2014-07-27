package config

import (
	//"encoding/json"
	//"io/ioutil"
	"os"
)

type CLIConfig struct {
	URL string `json:"url"`
}

func Load() (CLIConfig, error) {

	var url = os.Getenv("REKA_HOST")

	if url == "" {
		println("missing config environment variable: REKAHOST")
		os.Exit(1)
	}

	var config = new(CLIConfig)
	config.URL = url
	return *config, nil

	/*
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
	*/
}
