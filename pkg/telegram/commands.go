package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var startMessage = `Привет! Я - HR бот.Сейчас задам тебе несколько вопросов. Если готов(а) продолжить общение со мной, жми "Да".
Если выберешь "Да", соглашаешься с условиями обрабтки данных.Продолжим диалог? `

// var banMessage = "Вы заблокировали бота @PMIIHrBot"
// var warningBanMessage = "Вы не можете заблокировать других ботов"
var noQuestionMessage = `Хорошо, понял Вас! Пожалуйста,поделитесь со мной, что явялется причиной вашего отказа? Это поможет мне при последующем отборе
кандидатов.`
var startDialogMessage = `Отлично!Я очень рад!Тогда начнем наш диалог) `
var educationQuestion = `Скажи,есть ли у тебя высшее техническое образование?`
var choiseProfil = `Здравствуйте! На данный момент у нас открыт набор на следующие позиции:`

var choiseProfilKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации об образовании
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Go-разработчик", "Golang backend - developer"),
		tgbotapi.NewInlineKeyboardButtonData("Java-разработчик", "jun java dev"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Frontend-разработчик", "middle js dev"),
		tgbotapi.NewInlineKeyboardButtonData("Специалист DS", "middle data science"),
	),
)
var choisePosition = `Здравствуйте! На данный момент у нас открыт набор на следующие позиции:`

var choisePositionKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации об образовании
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Junior", "Junior"),
		tgbotapi.NewInlineKeyboardButtonData("Middle", "Middle"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Стажёр", "Intern"),
	),
)
var educationKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации об образовании
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "Have high technical education"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Нет", "Haven't high technical education"),
	),
)
var answerKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // inline меню для начала общения
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
var noQuestionKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о причинах отказа общаться с ботом
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
