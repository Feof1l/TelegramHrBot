package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var StartMessage = `Привет! Я - HR бот.Сейчас задам тебе несколько вопросов. Если готов(а) продолжить общение со мной, жми "Да".
Если выберешь "Да", соглашаешься с условиями обрабтки данных.Продолжим диалог? `
var BanMessage = "Вы заблокировали бота @PMIIHrBot"
var WarningBanMessage = "Вы не можете заблокировать других ботов"
var NoQuestionMessage = `Хорошо, понял Вас! Пожалуйста,поделитесь со мной, что явялется причиной вашего отказа? Это поможет мне при последующем отборе
кандидатов.`

var AnswerKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // inline меню для начала общения
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "Yes"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "No"),
		tgbotapi.NewInlineKeyboardButtonData("Заблокировать", "Block"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Пользовательское соглашение", UserAgreement),
	),
)
var NoQuestionKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о причинах отказа общаться с ботом
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вакансия неинтересна", "Вакансия неинтересна"),
		tgbotapi.NewInlineKeyboardButtonData("Уже нашел работу", "Уже нашел работу"),
		tgbotapi.NewInlineKeyboardButtonData("Другая причина", "Другая причина"),
		tgbotapi.NewInlineKeyboardButtonData("Не хочу говорить", "Не хочу говорить"),
	),
)
var BlockedUsers = make(map[string]bool)          // хэш мапа для хранения информации о заблокированных пользователях
var UserAgreement = "https://telegram.org/tos/ru" // сылка на пользовательское соглашение

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

	msgToDelete := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: msgID,
	}

	_, err := b.bot.MakeRequest(msgToDelete)

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
				b.errorLog.Println(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			switch update.CallbackQuery.Data {
			case "Block":
				BlockedUsers[b.bot.Self.UserName] = true
				err := b.clearChatHistory(update.Message.Chat.ID)
				if err != nil {
					b.errorLog.Println(err)
				}
			case "No":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, NoQuestionMessage)
				msg.ReplyMarkup = NoQuestionKeyBoard
				if _, err := b.bot.Send(msg); err != nil {
					b.errorLog.Println(err)
				}
				if update.CallbackQuery != nil {
					//логика добавления update.CallbackQuery.Data в БД
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо за обратную связь! Удачи!")
			case "Yes":

			default:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пожалуйста,используйте кнопки для общения с ботом")
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
func (b *Bot) IsBlockedUser() bool {
	for key, _ := range BlockedUsers {
		if key == b.bot.Self.UserName && BlockedUsers[key] == true {
			return false
		}
	}
	return true
}
