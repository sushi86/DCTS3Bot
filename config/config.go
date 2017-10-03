package config

import (
	"encoding/json"
	"os"
)

var config Configuration

type Configuration struct {
	DiscordApiKey	string
	TelegramApiKey	string
	TS3Hostname		string
	TS3Password		string
}

func readFile() {
	file, err := os.Open("conf.json")
	if err != nil {
		panic("no valid config file!")
	}
	decoder := json.NewDecoder(file)

	configuration := Configuration{}
	err = decoder.Decode(&configuration)

	if err != nil {
		panic("no valid json config file! Error: " + err.Error())
	}

	config = configuration
}

func GetConfig() Configuration {
	if config.DiscordApiKey == "" {
		readFile()
	}

	return config
}
