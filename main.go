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
)

// AppConfig : Init structure for config file.
type AppConfig struct {
	DeployToken  string
	GoogleSheets GSheets
	Telegram     TgOpts
	// AllowedUsers []string `yaml:allowed_users`
}

// GSheets : Sheets struct for Config File.
type GSheets struct {
	Token             string
	ID                string
	APISecretFileName string
}

// TgOpts : Telegram Options struct.
type TgOpts struct {
	Token string
	BotID string
}

// PrintConfig : Self-explanatory
func PrintConfig(Credentials *AppConfig) {
	fmt.Printf("%+v\n", Credentials)
}

// ReadConfig : Self-explanatory
func ReadConfig() (*AppConfig, error) {
	var Credentials = new(AppConfig)

	Credentials.GoogleSheets.Token = os.Getenv("GSHEETS_TOKEN")
	Credentials.GoogleSheets.ID = os.Getenv("GSHEETS_ID")
	Credentials.GoogleSheets.APISecretFileName = "gsheets_api_secret.json"

	Credentials.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	Credentials.Telegram.BotID = os.Getenv("TELEGRAM_BOT_ID")
	Credentials.DeployToken = os.Getenv("DEPLOY_TOKEN")

	return Credentials, nil
}

func telegramAnswerToUser() {
	// TODO: Описать входящие и возвратные аргументы

	Credentials := ReadConfig()
	token := Credentials.Telegram.Token

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
