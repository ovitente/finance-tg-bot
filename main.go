package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	GoogleSheets string       `yaml:google_sheets_token`
	Telegram     TelegramOpts `yaml:telegram`
	// AllowedUsers []string `yaml:allowed_users`
}

type TelegramOpts struct {
	Token string `yaml:token`
	BotId string `yaml:bot_id`
}

func appConf() *ConfigFile {
	var confPath = "./config.yaml"
	var cf = new(ConfigFile)

	file, err := os.Open(confPath)
	if err != nil {
		log.Fatal(err)
	}

	yamlByte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if err = yaml.Unmarshal(yamlByte, cf); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cf)
	return cf

}

func main() {
	appConf()
}
