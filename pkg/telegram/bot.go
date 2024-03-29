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
var StartDialogMessage = `Отлично!Я очень рад!Тогда начнем наш диалог) `
var EducationQuestion = `Скажи,есть ли у тебя высшее техническое образование?`
var EducationKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации об образовании
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "Have high technical education"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Нет", "Haven't high technical education"),
	),
)
var AnswerKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // inline меню для начала общения
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "Yes"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "No"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Заблокировать", "Block"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Пользовательское соглашение", UserAgreement),
	),
)
var NoQuestionKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о причинах отказа общаться с ботом
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вакансия неинтересна", "vacancy is not interesting"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Уже нашел работу", "already found a job"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Другая причина", "another reason"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Не хочу говорить", "don't want to talk"),
	),
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
	for update := range updates {

		if update.Message != nil && b.IsBlockedUser() { // потом переделать
			MessageIdDic[update.Message.MessageID]++
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
			b.SendMsg(msg)

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
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, NoQuestionMessage)

				msg.ReplyMarkup = NoQuestionKeyBoard
				b.SendMsg(msg)

			case "vacancy is not interesting":
				b.feedback(update.CallbackQuery)
			case "already found a job":
				b.feedback(update.CallbackQuery)
			case "another reason":
				b.feedback(update.CallbackQuery)
			case "don't want to talk":
				b.feedback(update.CallbackQuery)
			case "Yes":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, StartDialogMessage)
				b.SendMsg(msg)
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, EducationQuestion)
				msg.ReplyMarkup = EducationKeyBoard
				b.SendMsg(msg)
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
