package telegram

import (
	"log"

	"github.com/Feof1l/TelegramHrBot/pkg/models"
	"github.com/Feof1l/TelegramHrBot/pkg/models/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var BlockedUsers = make(map[string]bool) // хэш мапа для хранения информации о заблокированных пользователях
var MessageIdDic = make(map[int]int)

type Bot struct {
	bot        *tgbotapi.BotAPI
	errorLog   *log.Logger
	infoLog    *log.Logger
	candidates *mysql.CandidatModel
}

func NewBot(bot *tgbotapi.BotAPI, errorLog *log.Logger, infoLog *log.Logger, candidates *mysql.CandidatModel) *Bot {

	return &Bot{bot: bot, errorLog: errorLog, infoLog: infoLog, candidates: &mysql.CandidatModel{DB: candidates.DB}}
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
	queryCandidat := models.Possible_candidate{}
	queryPosition := models.Position{}
	flagNameCandidate := false
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
				b.SendMsg(msg)

			}
			if flagNameCandidate {
				queryCandidat.Candidate_name = update.Message.Text
				queryCandidat.Telegram_username = update.Message.Chat.UserName
				/*
					b.infoLog.Println(candidateName, update.Message.Chat.UserName)
					err := b.candidates.Insert(candidateName, update.Message.Chat.UserName)
					if err != nil {
						b.errorLog.Println(err)
					}*/
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приятно познакомиться!")
				b.SendMsg(msg)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберети специализацию, на которой хотите работать!")
				b.SendMsg(msg)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, choiseProfil)

				msg.ReplyMarkup = choiseProfilKeyBoard

				b.SendMsg(msg)

			}
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
				flagNameCandidate = true
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отлично! Для дальнейшего общения напишите, пожалуйста, мне свою фамилию и имя")
				b.SendMsg(msg)

			/*case "Yes":

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберети специализацию, на которой хотите работать!")
			b.SendMsg(msg)
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, choiseProfil)

			msg.ReplyMarkup = choiseProfilKeyBoard

			b.SendMsg(msg)*/
			case "Golang backend - developer", "Java backend - developer":
				queryPosition.Profil = update.CallbackQuery.Data
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберети позицию, на которой хотите работать!")
				b.SendMsg(msg)
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, choisePosition)
				msg.ReplyMarkup = choisePositionKeyBoard

				b.SendMsg(msg)
			case "Junior", "Middle", "Intern", "Senior", "Team Lead":
				queryPosition.Position_name = update.CallbackQuery.Data
				queryCandidat.Id_pos = DetermineId_pos(queryPosition.Profil, queryPosition.Position_name)
				err := b.candidates.Insert(queryCandidat.Candidate_name, queryCandidat.Telegram_username, queryCandidat.Id_pos)
				if err != nil {
					b.errorLog.Println(err)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, educationQuestion)
				msg.ReplyMarkup = educationKeyBoard
				b.SendMsg(msg)
			case "Среднее", "Неоконченное высшее", "Бакалавр", "Магистр", "Кандидат наук":
				queryCandidat.Education = update.CallbackQuery.Data

				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.Update(queryCandidat.Education, queryCandidat.Candidate_name)
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.CallCompareEducation(queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)

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
