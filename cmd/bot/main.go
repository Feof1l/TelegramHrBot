package main

import (
	"log"
	"os"

	"github.com/Feof1l/TelegramHrBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	// создаем новый логер для вывода информационных сообщенйи в поток stdout c припиской info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//аналогично для логов с ошибками, такеж включим вывод фйла и номера  строки, где произошла ошибка
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	/*app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}*/
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		errorLog.Panic(err)
	}

	bot.Debug = true

	infoLog.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		// универсальный ответ на любое сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		reply := ""
		if update.Message == nil { // If we got a message
			//infoLog.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			continue
		}

		switch update.Message.Command() {
		case "start":
			reply = telegram.StartMessage
			msg.ReplyMarkup = telegram.AnswerKeyBoard
		}
		switch update.Message.Text {
		case "open":
			//msg.ReplyMarkup = numericKeyboard
		case "close":
			//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		}
		// создаем ответное сообщение
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		if _, err := bot.Send(msg); err != nil {
			errorLog.Panic(err)
		}

	}
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}
