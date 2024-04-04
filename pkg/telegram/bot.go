package telegram

import (
	"log"

	"github.com/Feof1l/TelegramHrBot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var BlockedUsers = make(map[string]bool)          // хэш мапа для хранения информации о заблокированных пользователях
var UserAgreement = "https://telegram.org/tos/ru" // сылка на пользовательское соглашение
var MessageIdDic = make(map[int]int)

type Bot struct {
	bot      *tgbotapi.BotAPI
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewBot(bot *tgbotapi.BotAPI, errorLog *log.Logger, infoLog *log.Logger) *Bot {

	return &Bot{bot: bot, errorLog: errorLog, infoLog: infoLog}
}
func (b *Bot) Start() error {

	b.infoLog.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)

	return nil
}
func (b *Bot) clearChatHistory(chatID int64) error {
	for key := range MessageIdDic {
		msgToDelete := tgbotapi.DeleteMessageConfig{
			ChatID:    chatID,
			MessageID: key,
		}

		_, err := b.bot.DeleteMessage(msgToDelete)
		if err != nil {
			b.errorLog.Println(err)
			return err
		}
	}

	return nil
}
func (b *Bot) SendMsg(msg tgbotapi.MessageConfig) error {
	sendMessage, err := b.bot.Send(msg)
	if err != nil {
		b.errorLog.Println(err)
		return err
	}
	MessageIdDic[sendMessage.MessageID]++
	return nil
}
func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	queryPosition := models.Position{}
	for update := range updates {

		if update.Message != nil && b.IsBlockedUser() { // потом переделать
			MessageIdDic[update.Message.MessageID]++
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Command() {
			case "start":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, startMessage)
				msg.ReplyMarkup = answerKeyBoard

			}
			b.SendMsg(msg)

			// Send the message.

		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			switch update.CallbackQuery.Data {
			case "Block":
				BlockedUsers[b.bot.Self.UserName] = true
				err := b.clearChatHistory(update.CallbackQuery.Message.Chat.ID)
				if err != nil {
					b.errorLog.Println(err)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы заблокировали бота, спасибо за общение!")
				b.SendMsg(msg)
			case "No":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, noQuestionMessage)

				msg.ReplyMarkup = noQuestionKeyBoard
				b.SendMsg(msg)

			case "vacancy is not interesting", "another reason", "don't want to talk", "already found a job":
				b.feedback(update.CallbackQuery)
			case "Yes":
				/*msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, StartDialogMessage)
				b.SendMsg(msg)
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, EducationQuestion)
				msg.ReplyMarkup = EducationKeyBoard
				b.SendMsg(msg)*/
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберети специализацию, на которой хотите работать!")
				b.SendMsg(msg)
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, choiseProfil)

				msg.ReplyMarkup = choiseProfilKeyBoard

				b.SendMsg(msg)
			case "Golang backend - developer", "Java backend - developer":
				queryPosition.Profil = update.CallbackQuery.Data
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберети позицию, на которой хотите работать!")
				b.SendMsg(msg)
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, choisePosition)
				msg.ReplyMarkup = choisePositionKeyBoard

				b.SendMsg(msg)
			case "Junior", "Middle", "Intern":
				queryPosition.Position_name = update.CallbackQuery.Data

			default:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пожалуйста,используйте кнопки для общения с ботом")
				b.SendMsg(msg)

			}

		}
	}

}
func (b *Bot) feedback(CallbackQuery *tgbotapi.CallbackQuery) {
	//логика добавления update.CallbackQuery.Data в БД
	msg := tgbotapi.NewMessage(CallbackQuery.Message.Chat.ID, "Спасибо за обратную связь! Удачи!")
	b.SendMsg(msg)

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
