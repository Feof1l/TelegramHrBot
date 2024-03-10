package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

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
		tgbotapi.NewInlineKeyboardButtonURL("Пользовательское соглашение", "http://1.com"),
	),
)
