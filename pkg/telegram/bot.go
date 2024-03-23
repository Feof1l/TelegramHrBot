package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var StartMessage = `Привет! Я - HR бот.Сейчас задам тебе несколько вопросов. Если готов(а) продолжить общение со мной, жми "Да".
Если выберешь "Да", соглашаешься с условиями обрабтки данных.Продолжим диалог? `
var BanMessage = "Вы заблокировали бота @PMIIHrBot"
var WarningBanMessage = "Вы не можете заблокировать других ботов"
var AnswerKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "Да"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "Нет"),
		tgbotapi.NewInlineKeyboardButtonData("Заблокировать", "Заблокировать"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Пользовательское соглашение", UserAgreement),
	),
)
var BlockedUsers = make(map[int64]bool)
var UserAgreement = "https://telegram.org/tos/ru"

type Bot struct {
	bot      *tgbotapi.BotAPI
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewBot(bot *tgbotapi.BotAPI, errorLog *log.Logger, infoLog *log.Logger) *Bot {

	return &Bot{bot: bot, errorLog: errorLog, infoLog: infoLog}
}
func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)

	return nil
}
func (b *Bot) clearChatHistory(chatID int64) error {
	// Получаем последнее обновление в чате
	updates, err := b.bot.GetUpdates(tgbotapi.UpdateConfig{Timeout: 1})
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
		_, err := b.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    chatID,
			MessageID: i,
		})
		if err != nil {
			b.errorLog.Println("Failed to delete message:", err)
		}
	}

	return nil
}
func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, StartMessage)
				msg.ReplyMarkup = AnswerKeyBoard

			}

			// Send the message.
			if _, err := b.bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			switch update.CallbackQuery.Data {
			case "Заблокировать":
				err := b.clearChatHistory(update.Message.Chat.ID)
				if err != nil {
					b.errorLog.Println(err)
				}
				/*msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}*/

			default:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "AAAAAAAAAAAAAAAAAA")
				if _, err := b.bot.Send(msg); err != nil {
					b.errorLog.Println(err)
				}

			}

		}
	}

}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
