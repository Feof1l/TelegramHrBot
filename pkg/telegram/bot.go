package telegram

import (
	"log"
	"strconv"

	"github.com/Feof1l/TelegramHrBot/pkg/models"
	"github.com/Feof1l/TelegramHrBot/pkg/models/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var BlockedUsers = make(map[string]bool) // хэш мапа для хранения информации о заблокированных пользователях
var MessageIdDic = make(map[int]int)
var flagFeedback = false

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

// метод удаления всех сообщений
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

// обёртка для отправки сообщения ботом
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

	flagFeadbackAnotherReason := false
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
			if flagFeedback {
				if err := b.candidates.InsertFeadBack(update.Message.Text, queryCandidat.Id_possible_candidate); err != nil {
					b.errorLog.Println(err)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо за обратную связь! Удачи!")
				b.SendMsg(msg)
				flagFeedback = !flagFeedback
			}
			if flagFeadbackAnotherReason {

				err := b.candidates.Insert(update.Message.Chat.UserName, update.Message.Chat.UserName, 1)
				if err != nil {
					b.errorLog.Println(err)
				}
				id, err := b.candidates.GetId(update.Message.Chat.UserName, update.Message.Chat.UserName)
				if err != nil {
					b.errorLog.Println(err)
				}
				if err = b.candidates.InsertFeadBack(update.Message.Text, id); err != nil {
					b.errorLog.Println(err)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо за обратную связь! Удачи!")
				b.SendMsg(msg)
				flagFeadbackAnotherReason = !flagFeadbackAnotherReason

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
				flagNameCandidate = !flagNameCandidate

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
			case "Вакансия неинтересна", "Не хочу говорить", "Уже нашел работу":

				err := b.candidates.Insert(update.CallbackQuery.Message.Chat.UserName, update.CallbackQuery.Message.Chat.UserName, 1)
				if err != nil {
					b.errorLog.Println(err)
				}
				id, err := b.candidates.GetId(update.CallbackQuery.Message.Chat.UserName, update.CallbackQuery.Message.Chat.UserName)
				if err != nil {
					b.errorLog.Println(err)
				}
				if err := b.candidates.InsertFeadBack(update.CallbackQuery.Data, id); err != nil {
					b.errorLog.Println(err)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Спасибо за обратную связь! Удачи!")
				b.SendMsg(msg)

			case "Другая причина":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пожалуйста,напишите мне, почему Вы не хотите со мной общаться")
				b.SendMsg(msg)
				flagFeadbackAnotherReason = true

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
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите позицию, на которой хотите работать!")
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
				err = b.candidates.UpdateStringData("Education", queryCandidat.Education, queryCandidat.Id_possible_candidate)
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.CallStoredProcedure("Compare_Education", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFailFlag(id)
				if err != nil {
					b.errorLog.Println(err)
				}
				//b.infoLog.Println(failFlag)
				if failFlag {
					b.failFlag(update)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, citizenshipMessage)
					msg.ReplyMarkup = citizenshipKeyBoard
					b.SendMsg(msg)
				}

			case "РФ", "РБ", "СНГ", "Другое":
				queryCandidat.Citizenship = update.CallbackQuery.Data
				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.UpdateStringData("Citizenship", queryCandidat.Citizenship, queryCandidat.Id_possible_candidate)
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.CallStoredProcedure("Compare_Citizenship", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFailFlag(id)
				if err != nil {
					b.errorLog.Println(err)
				}

				if failFlag {
					b.failFlag(update)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, workExperienceMessage)
					msg.ReplyMarkup = workExperienceKeyBoard
					b.SendMsg(msg)
				}

			case "Менее года", "1 - 3 года", "3 - 6 лет", "Более 6 лет":
				queryCandidat.Work_experience = update.CallbackQuery.Data
				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.UpdateStringData("Work_experience", queryCandidat.Work_experience, queryCandidat.Id_possible_candidate)
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.CallStoredProcedure("Compare_Work_experience", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFailFlag(id)
				if err != nil {
					b.errorLog.Println(err)
				}

				if failFlag {
					b.failFlag(update)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, hoursMessage)
					msg.ReplyMarkup = hoursKeyBoard
					b.SendMsg(msg)
				}

			case "20", "30", "40":
				hours, err := strconv.Atoi(update.CallbackQuery.Data)
				if err != nil {
					b.errorLog.Println(err)
				}
				queryCandidat.Hours = hours
				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.UpdateIntData("Hours", queryCandidat.Hours, queryCandidat.Id_possible_candidate)
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.CallStoredProcedure("Compare_hours", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFailFlag(id)
				if err != nil {
					b.errorLog.Println(err)
				}

				if failFlag {
					b.failFlag(update)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, workFormatMessage)
					msg.ReplyMarkup = workFormatKeyBoard
					b.SendMsg(msg)
				}

			default:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пожалуйста,используйте кнопки для общения с ботом")
				b.SendMsg(msg)

			}

		}
	}

}
func (b *Bot) failFlag(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, failTrueMessage)
	b.SendMsg(msg)
	msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, feedbackMessage)
	b.SendMsg(msg)
	flagFeedback = true
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
