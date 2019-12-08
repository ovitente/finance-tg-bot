package main

import (
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

// ConfigFile : Init structure for config file.
type ConfigFile struct {
	GoogleSheets GSheets `yaml:"google_sheets"`
	Telegram     TgOpts  `yaml:"telegram"`
	// AllowedUsers []string `yaml:allowed_users`
}

// GSheets : Sheets struct for Config File.
type GSheets struct {
	Token             string `yaml:"token,omitempty"`
	ID                string `yaml:"id,omitempty"`
	APISecretFileName string `yaml:"api_secret_file_name,omitempty"`
}

// TgOpts : Telegram Options struct.
type TgOpts struct {
	Token string `yaml:"token"`
	BotID string `yaml:"bot_id"`
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func appConf() *ConfigFile {

	var configFilePath = "./config.yml"
	var configFile = new(ConfigFile)

	file, err := os.Open(configFilePath)
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

	return configFile
}

func telegramAnswerToUser() {
	// TODO: Описать входящие и возвратные аргументы

	configFile := appConf()
	token := configFile.Telegram.Token

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

func main() {
	// Read config
	configFile := appConf()

	// GOOGLE SHEETS
	data, err := ioutil.ReadFile(configFile.GoogleSheets.APISecretFileName)
	checkError(err)

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(configFile.GoogleSheets.ID)
	checkError(err)

	sheet, err := spreadsheet.SheetByIndex(0)
	checkError(err)

	// Printing content of the whole table

	// for _, row := range sheet.Rows {
	// 	for _, cell := range row {
	// 		fmt.Println(cell.Value)
	// 	}
	// }

	// Update cell content
	sheet.Update(0, 0, "CELL EDITED WITH GOLANG")

	// Make sure call Synchronize to reflect the changes
	err = sheet.Synchronize()
	checkError(err)

	// Init telegram bot
	telegramAnswerToUser()
}
