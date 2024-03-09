package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var StartMessage = "Привет! Я - HR бот.Сейчас задам тебе несколько вопросов. Если готов(а) продолжить общение со мной, жми `ДА`. \nЕсли выберешь `ДА`, соглашаешься с условиями обрабтки данных.Продолжим диалог? "

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

var AnswerKeyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Да"),
		tgbotapi.NewKeyboardButton("Нет"),
		tgbotapi.NewKeyboardButton("Заблокировать"),
	),
)
