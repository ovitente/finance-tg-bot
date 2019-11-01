package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
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

func telegramSendBack() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("ТОКЕН")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	err = bot.UpdatesChan(ucfg)
	// читаем обновления из канала
	for {
		select {
		case update := <-bot.Updates:
			// Пользователь, который написал боту
			UserName := update.Message.From.UserName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			ChatID := update.Message.Chat.ID

			// Текст сообщения
			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			// Ответим пользователю его же сообщением
			reply := Text
			// Созадаем сообщение
			msg := tgbotapi.NewMessage(ChatID, reply)
			// и отправляем его
			bot.SendMessage(msg)
		}

	}
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

}

func main() {
	appConf()
	telegramSendBack()
