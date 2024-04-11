package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var UserAgreement = "https://telegram.org/tos/ru" // сылка на пользовательское соглашение
var startMessage = `Привет! Я - HR бот.Сейчас задам тебе несколько вопросов. Если готов(а) продолжить общение со мной, жми "Да".
Если выберешь "Да", соглашаешься с условиями обрабтки данных.Продолжим диалог? `

// var banMessage = "Вы заблокировали бота @PMIIHrBot"
// var warningBanMessage = "Вы не можете заблокировать других ботов"
var noQuestionMessage = `Хорошо, понял тебя! Пожалуйста,поделись со мной, что явялется причиной твоего отказа? Это поможет мне при последующем отборе
кандидатов.`
var salaryKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о зарплате
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да,все верно", "Correct salary"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Нет(изменить)", "Uncorrect salary"),
	),
)
var startDialogMessage = `Отлично! Я очень рад! Тогда начнем наш диалог) `
var educationQuestion = `Какое у тебя образование. Выбери подходящий вариант?`
var choiseProfil = `Здравствуйте! На данный момент нам требуются следующие специалисты`
var failTrueMessage = `Ой, к сожалению, это не подходит для данной вакансии, и мы с не сможем рассмотреть 
твою кандидатуру 😢. Но не расстраивайся, у нас в компании есть много других отличных вакансий,
Будем рады, если ты найдешь что-то подходящее!`
var feedbackMessage = `Знаешь, мне было приятно с тобой пообщаться. Поэтому я сохраню твой контакт на 
случай появления подходящих вакансий в компании.
Надеюсь, и тебе было полезно со мной поговорить. Скажи, а что тебе понравилось при взаимодействии со мной?
Напиши ответ в свободной форме`
var expectedSalaryMessage = `Я тебя понял! На какую зарплаты ты расчитываешь?`
var workFormatMessage = `Нам подходит! А какой формат работы ты хочешь?`
var workFormatKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о формате работы
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Оффис(полностью оффлайн)", "Оффис"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удаленка", "Удаленка"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Гибрид(оффис + удаленка)", "Гибрид"),
	),
)
var hoursMessage = `Отлично! Скажи, какую занятость ты рассматриваешь?`
var hoursKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о занятости
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("20 часов в неделю", "20"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("30 часов в неделю", "30"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("40 часов в неделю", "40"),
	),
)
var workExperienceMessage = `Класс! Расскажи, какой у тебя опыт коммерческой разработки,сколько лет работаешь?`
var workExperienceKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации об опыте разработки
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Менее года", "Менее года"),
		tgbotapi.NewInlineKeyboardButtonData("От 1 до 3 лет", "1 - 3 года"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("От 3 до 6 лет", "3 - 6 лет"),
		tgbotapi.NewInlineKeyboardButtonData("Более 6 лет", "Более 6 лет"),
	),
)
var citizenshipMessage = `Супер!Скажи, какое у тебя гражданство?`
var citizenshipKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о гражданстве
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("РФ", "РФ"),
		tgbotapi.NewInlineKeyboardButtonData("РБ", "РБ"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("СНГ", "СНГ"),
		tgbotapi.NewInlineKeyboardButtonData("Другое", "Другое"),
	),
)
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
var choisePosition = `На данный момент у нас открыт набор на следующие позиции:`

var choisePositionKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации об образовании
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Стажёр", "Intern"),
		tgbotapi.NewInlineKeyboardButtonData("Junior", "Junior"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Middle", "Middle"),
		tgbotapi.NewInlineKeyboardButtonData("Senior", "Senior"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Team Lead", "Team Lead"),
	),
)
var educationKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации об образовании
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Среднее", "Среднее"),
		tgbotapi.NewInlineKeyboardButtonData("Неоконченное высшее", "Неоконченное высшее"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Высшее (бакалавриат)", "Бакалавр"),
		tgbotapi.NewInlineKeyboardButtonData("Высшее (магистратура)", "Магистр"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Высшее (аспирантура)", "Кандидат наук"),
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
		tgbotapi.NewInlineKeyboardButtonData("Вакансия неинтересна", "Вакансия неинтересна"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Уже нашел работу", "Уже нашел работу"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Другая причина", "Другая причина"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Не хочу говорить", "Не хочу говорить"),
	),
)

func DetermineId_pos(profil, position string) int {
	if profil == "Golang backend - developer" && position == "Junior" {
		return 1
	} else if profil == "java backend - developer" && position == "Middle" {
		return 2
	} else {
		return 3
	}
}
