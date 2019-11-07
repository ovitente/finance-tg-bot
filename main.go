package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

// ConfigFile Init structure for config file.
type ConfigFile struct {
	GoogleSheets GSheets `yaml:"google_sheets"`
	Telegram     TgOpts  `yaml:"telegram"`
	// AllowedUsers []string `yaml:allowed_users`
}

// Sheets struct.
type GSheets struct {
	Token string `yaml:"token,omitempty"`
	ID    string `yaml:"id,omitempty"`
}

// Telegram Options struct.
type TgOpts struct {
	Token string `yaml:"token"`
	BotId string `yaml:"bot_id"`
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}

	// TELEGRAM
	token := os.Getenv("TOKEN")

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		URL:    "https://api.telegram.org",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()
}

func appConf() *ConfigFile {

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

	// TODO: refactor return
	return configFile

}

func main() {
	appConf()
	// GOOGLE SHEETS
	data, err := ioutil.ReadFile("client_secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(config)
	checkError(err)
	sheet, err := spreadsheet.SheetByIndex(0)
	checkError(err)

	for _, row := range sheet.Rows {
		for _, cell := range row {
			fmt.Println(cell.Value)
		}
	}

	// Update cell content
	sheet.Update(0, 0, "hogehoge")

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(err)
}
