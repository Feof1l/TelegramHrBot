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
		log.Fatal(err)
	}
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		errorLog.Panic(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, errorLog, infoLog)
	func() {
		if err := telegramBot.Start(); err != nil {
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
func clearChatHistory(bot *tgbotapi.BotAPI, chatID int64) error {
	// Получаем последнее обновление в чате
	updates, err := bot.GetUpdates(tgbotapi.UpdateConfig{Timeout: 1})
	if err != nil {
		return err
	}

	// Находим ID последнего сообщения
	var lastMessageID int
	for _, update := range updates {
		if update.Message != nil && update.Message.Chat.ID == chatID {
			lastMessageID = update.Message.MessageID
		}
	}

	// Удаляем сообщения в диапазоне от 1 до последнего сообщения
	for i := 1; i <= lastMessageID; i++ {
		_, err := bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    chatID,
			MessageID: i,
		})
		if err != nil {
			log.Println("Failed to delete message:", err)
		}
	}

	return nil
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
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

/*
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
			// Обработка нажатия на кнопку во всплывающем окне
			if update.CallbackQuery.Data == "button_pressed" {
				// Здесь можно добавить логику обработки нажатия на кнопку
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы нажали на кнопку!")
				bot.Send(msg)
			}
		} else if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Text {
			case "/start":
				// Создаем всплывающую клавиатуру
				inlineBtn := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Нажми меня", "button_pressed"),
					),
				)

				msg.Text = "Привет! Нажми на кнопку, чтобы получить сообщение."
				msg.ReplyMarkup = inlineBtn
			default:
				msg.Text = "Я не понимаю, что вы имеете в виду."
			}
		}
		switch update.Message.Text {

		case "close":
			//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		default:
			//continue
		}
		// создаем ответное сообщение
*/
