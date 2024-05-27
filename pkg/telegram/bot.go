package telegram

import (
	"errors"
	"log"
	"strconv"
	"time"

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

var ErrUncorrectSalary = errors.New("Зарплата указана некорректно")
var ErrUncorrectContactNumber = errors.New("Контактный номер указан некорректно!")

func NewBot(bot *tgbotapi.BotAPI, errorLog *log.Logger, infoLog *log.Logger, candidates *mysql.CandidatModel) *Bot {

	return &Bot{bot: bot, errorLog: errorLog, infoLog: infoLog, candidates: &mysql.CandidatModel{DB: candidates.DB}}
}
func (b *Bot) Start() error {

	b.infoLog.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.HandleUpdates(updates)

	return nil
}

// метод удаления всех сообщений
func (b *Bot) ClearChatHistory(chatID int64) error {
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
func (b *Bot) HandleUpdates(updates tgbotapi.UpdatesChannel) {
	queryCandidat := models.Possible_candidate{}
	queryPosition := models.Position{}
	flagNameCandidate := false
	flagExpectedSalary := false
	flagContactNumber := false

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
			switch {
			case flagFeedback:
				if err := b.candidates.InsertFeadBack(update.Message.Text, queryCandidat.Id_possible_candidate); err != nil {
					b.errorLog.Println(err)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо за обратную связь! Удачи!")
				b.SendMsg(msg)
				flagFeedback = !flagFeedback
			case flagFeadbackAnotherReason:
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
			case flagNameCandidate:
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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, choiseProfilMessage)
				b.SendMsg(msg)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, choiseProfil)

				msg.ReplyMarkup = choiseProfilKeyBoard

				b.SendMsg(msg)
				flagNameCandidate = !flagNameCandidate
			case flagExpectedSalary:
				expectedSalaryString := update.Message.Text
				expectedSalaryInt, err := strconv.Atoi(expectedSalaryString)
				queryCandidat.Expected_salary = expectedSalaryInt
				if err != nil || queryCandidat.Expected_salary <= 0 {
					b.errorLog.Println(ErrUncorrectSalary)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Зарплата указана некорректно!Пожалуйста, введите число")
					b.SendMsg(msg)

				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы рассчитываете на зарплату:"+strconv.Itoa(queryCandidat.Expected_salary)+", так?")
					msg.ReplyMarkup = salaryKeyBoard
					b.SendMsg(msg)
					flagExpectedSalary = !flagExpectedSalary
				}
			case flagContactNumber:
				contactNumberString := update.Message.Text
				queryCandidat.Contact_number = contactNumberString
				if !isValidPhoneNumber(contactNumberString) {
					b.errorLog.Println(ErrUncorrectContactNumber)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, `Номер указан некорректно!Пожалуйста, введите номер в одном из форматов:"8xxxxxxxxxx", "+7xxxxxxxxxx", "8 (xxx) xxx-xx-xx", "+7 xxx xxx-xx-xx" `)
					b.SendMsg(msg)
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Вы указали свой номер телефона:"+queryCandidat.Contact_number+", так?")
					msg.ReplyMarkup = contactNumberKeyBoard
					b.SendMsg(msg)
					flagContactNumber = !flagContactNumber
				}
			}

			// Send the message.

		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			switch update.CallbackQuery.Data {
			case "Block":
				BlockedUsers[b.bot.Self.UserName] = true
				err := b.ClearChatHistory(update.CallbackQuery.Message.Chat.ID)
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
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, candidateNameMessage)
				b.SendMsg(msg)

			/*case "Yes":

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберети специализацию, на которой хотите работать!")
			b.SendMsg(msg)
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, choiseProfil)

			msg.ReplyMarkup = choiseProfilKeyBoard

			b.SendMsg(msg)*/
			case "Golang backend - developer", "Java backend - developer", "Frontend - developer", "Data Science - specialist":
				queryPosition.Profil = update.CallbackQuery.Data
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, choisePositionMessage)
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

				currentTime := time.Now()
				currentTimetring := currentTime.Format("02.01.2006")
				queryCandidat.Date_of_dialog = currentTimetring
				b.infoLog.Println(queryCandidat.Date_of_dialog)
				err = b.candidates.UpdateStringData("Date_of_dialog", queryCandidat.Date_of_dialog, queryCandidat.Id_possible_candidate)
				if err != nil {
					b.errorLog.Println(err)
				}

				err = b.candidates.UpdateStringData("Education", queryCandidat.Education, queryCandidat.Id_possible_candidate)
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.CallStoredProcedure("Compare_Education", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFlag("Fail_flag", id)
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
				failFlag, err := b.candidates.GetFlag("Fail_flag", id)
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
				failFlag, err := b.candidates.GetFlag("Fail_flag", id)
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
				failFlag, err := b.candidates.GetFlag("Fail_flag", id)
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
			case "Оффис", "Удаленка", "Гибрид":

				queryCandidat.Work_format = update.CallbackQuery.Data
				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.UpdateStringData("Work_format", queryCandidat.Work_format, queryCandidat.Id_possible_candidate)
				if err != nil {
					b.errorLog.Println(err)
				}
				err = b.candidates.CallStoredProcedure("Compare_work_format", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFlag("Fail_flag", id)
				if err != nil {
					b.errorLog.Println(err)
				}

				if failFlag {
					b.failFlag(update)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, expectedSalaryMessage)

					b.SendMsg(msg)
					flagExpectedSalary = true
				}
			case "Correct salary":
				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				b.candidates.UpdateIntData("Expected_salary", queryCandidat.Expected_salary, queryCandidat.Id_possible_candidate)
				err = b.candidates.CallStoredProcedure("Compare_Salary", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFlag("Fail_flag", id)
				if err != nil {
					b.errorLog.Println(err)
				}

				if failFlag {
					b.failFlag(update)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, readyToRelocateMessage)
					msg.ReplyMarkup = readyToRelocateKeyBoard
					b.SendMsg(msg)
				}
			case "Uncorrect salary":
				flagExpectedSalary = true
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, expectedSalaryMessage)
				b.SendMsg(msg)
			case "Ready to relocate", "Not ready to relocate":
				if update.CallbackQuery.Data == "Ready to relocate" {
					queryCandidat.Ready_to_relocate = true
				} else {
					queryCandidat.Ready_to_relocate = false
				}
				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				b.candidates.UpdateBoolData("Ready_to_relocate", queryCandidat.Ready_to_relocate, queryCandidat.Id_possible_candidate)
				err = b.candidates.CallStoredProcedure("Compare_relocation", queryCandidat.Id_pos, queryCandidat.Id_possible_candidate)
				failFlag, err := b.candidates.GetFlag("Fail_flag", id)
				if err != nil {
					b.errorLog.Println(err)
				}

				if failFlag {
					b.failFlag(update)
				} else {
					err := b.candidates.UpdateBoolData("ready_flag", true, id)
					if err != nil {
						b.errorLog.Println(err)
					} else {
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, goodResultMessage)

						b.SendMsg(msg)

						msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, contactNumberMessage)
						flagContactNumber = true
						b.SendMsg(msg)

					}

				}
			case "Correct number":
				id, err := b.candidates.GetId(queryCandidat.Candidate_name, queryCandidat.Telegram_username)
				queryCandidat.Id_possible_candidate = id
				if err != nil {
					b.errorLog.Println(err)
				}
				b.candidates.UpdateStringData("Contact_number", queryCandidat.Contact_number, queryCandidat.Id_possible_candidate)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, finalMessage)
				flagFeedback = true
				b.SendMsg(msg)

			case "Uncorrect number":
				flagContactNumber = true
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, contactNumberMessage)
				b.SendMsg(msg)

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
