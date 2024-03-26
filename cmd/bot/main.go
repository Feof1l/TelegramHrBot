package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Feof1l/TelegramHrBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config struct {
	TelegramBotToken string
}

func main() {

	// создаем новый логер для вывода информационных сообщенйи в поток stdout c припиской info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//аналогично для логов с ошибками, такеж включим вывод фйла и номера  строки, где произошла ошибка
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// декодируем файл json, в котором хранится конфиг - токен бота

	/*app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}*/
	configuration, err := decodeConfig("config.json")
	if err != nil {
		errorLog.Println(err)
	}
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		errorLog.Println(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, errorLog, infoLog)
	func() {
		if err := telegramBot.Start(); err != nil && !telegramBot.IsBlockedUser() {
			errorLog.Println(err)
		}
	}()

	////////////////////////

}
func decodeConfig(fileName string) (Config, error) {
	file, _ := os.Open(fileName)
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	return configuration, err
}

/*
			// универсальный ответ на любое сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//reply := ""
			if update.Message == nil { // If we got a message
				//infoLog.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				continue
			}

			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, telegram.StartMessage)
				msg.ReplyMarkup = telegram.AnswerKeyBoard

			}
			if update.CallbackQuery != nil {
				// Обработка нажатия на кнопку из inline меню
				if update.CallbackQuery.Data == "Заблокировать" {
					bot.Send(msg)
				}
			}
			switch update.Message.Text {
			case "open":
				//msg.ReplyMarkup = numericKeyboard
			case "close":
				//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			}
			// создаем ответное сообщение

			if _, err := bot.Send(msg); err != nil {
				errorLog.Println(err)
			}

		}
	}
*/
