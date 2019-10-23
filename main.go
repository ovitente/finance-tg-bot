package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	GoogleSheets GSheets `yaml:"google_sheets"`
	Telegram     TgOpts  `yaml:"telegram"`
	// AllowedUsers []string `yaml:allowed_users`
}

type GSheets struct {
	Token string `yaml:"token,omitempty"`
	ID    string `yaml:"id,omitempty"`
}
type TgOpts struct {
	Token string `yaml:"token"`
	BotId string `yaml:"bot_id"`
}

func main() {
	var confPath = "./config.yml"
	var configFile = new(ConfigFile)

	file, err := os.Open(confPath)
	if err != nil {
		log.Fatal(err)
	}

	yamlByte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if err = yaml.Unmarshal(yamlByte, configFile); err != nil {
		log.Fatal(err)
	}

	fmt.Println(configFile.GoogleSheets.Token)
	fmt.Println(configFile.GoogleSheets.ID)
	fmt.Println(configFile.Telegram.Token)
	fmt.Println(configFile.Telegram.BotId)

}
