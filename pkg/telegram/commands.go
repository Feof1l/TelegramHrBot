package telegram

import (
	"encoding/json"
	"os"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// сылка на пользовательское соглашение
var UserAgreement = "https://telegram.org/tos/ru"

// стартвоое сообщение, с которого начинается диалог
var startMessage = `Привет! Я - HR бот.Сейчас задам тебе несколько вопросов. Если готов(а) продолжить общение со мной, жми "Да".
Если выберешь "Да", соглашаешься с условиями обрабтки данных.Продолжим диалог? `

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

// сообщение с вопросом об имени кандидата
var candidateNameMessage = `Отлично! Для дальнейшего общения напишите, пожалуйста, мне свою фамилию и имя`

// сообщение с вопросом о специализации
var choiseProfilMessage = ` Выберети специализацию, на которой хотите работать!`

var choiseProfil = `Здравствуйте! На данный момент нам требуются следующие специалисты`

var choiseProfilKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о специализации
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Go-разработчик", "Golang backend - developer"),
		tgbotapi.NewInlineKeyboardButtonData("Java-разработчик", "Java backend - developer"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Frontend-разработчик", "Frontend - developer"),
		tgbotapi.NewInlineKeyboardButtonData("Специалист DS", "Data Science - specialist"),
	),
)

// сообщение с вопросом о позиции x
var choisePositionMessage = `Выберите позицию, на которой хотите работать!`

var choisePosition = `На данный момент у нас открыт набор на следующие позиции:`

var choisePositionKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о позиции
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

// сообщение с вопросом об образовании
var educationQuestion = `Какое у тебя образование. Выбери подходящий вариант?`

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

// сообщение с вопросом о гражданстве
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

// сообщение с вопросом об опыте работы
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

// сообщение с вопросом об уровне занятости
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

// сообщение с вопросом о формате работы
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

// сообщение с вопросом об уровне зарплаты
var expectedSalaryMessage = `Я тебя понял! На какую зарплаты ты расчитываешь?`

var salaryKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о зарплате
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да,все верно", "Correct salary"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Нет(изменить)", "Uncorrect salary"),
	),
)

// сообщение с вопросом о командировках
var readyToRelocateMessage = `Хорошо, а подскажи, готов ли ты к командировкам, переездам?`

var readyToRelocateKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о командировках
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Готов", "Ready to relocate"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Не готов", "Not ready to relocate"),
	),
)

// сообщение, когда кандидат успешно ответил на все вопросы
var goodResultMessage = `Благодарю за ответы!

Я передам наш разговор рекрутеру и он свяжется с тобой☎
`

var contactNumberMessage = `Для этого (на всякий случай) укажи,пожалуйста, свой номер телефона`

var contactNumberKeyBoard = tgbotapi.NewInlineKeyboardMarkup( // // inline меню для сборе инофрмации о номере телефона
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да,это мой номер", "Correct number"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Нет(изменить)", "Uncorrect number"),
	),
)

// var banMessage = "Вы заблокировали бота @PMIIHrBot"
// var warningBanMessage = "Вы не можете заблокировать других ботов"
var noQuestionMessage = `Хорошо, понял тебя! Пожалуйста,поделись со мной, что явялется причиной твоего отказа? Это поможет мне при последующем отборе
кандидатов.`

var startDialogMessage = `Отлично! Я очень рад! Тогда начнем наш диалог) `

var failTrueMessage = `Ой, к сожалению, это не подходит для данной вакансии, и мы с не сможем рассмотреть 
твою кандидатуру 😢. Но не расстраивайся, у нас в компании есть много других отличных вакансий,
Будем рады, если ты найдешь что-то подходящее!`
var feedbackMessage = `Знаешь, мне было приятно с тобой пообщаться. Поэтому я сохраню твой контакт на 
случай появления подходящих вакансий в компании.
Надеюсь, и тебе было полезно со мной поговорить. Скажи, а что тебе понравилось при взаимодействии со мной?
Напиши ответ в свободной форме`

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

var finalMessage = `Спасибо! На этом, к сожалению, наш разговор подошёл к концу! 
Мне было приятно с тобой общаться, надеюсь, тебе тоже! Если хочешь, можешь оставить свое впечатление о нашем диалоге, мне будет приятно!`

func DetermineId_pos(profil, position string) int {
	if profil == "Golang backend - developer" && position == "Junior" {
		return 1
	} else if profil == "Golang backend - developer" && position == "Intern" {
		return 2
	} else if profil == "Golang backend - developer" && position == "Middle" {
		return 3
	} else if profil == "Golang backend - developer" && position == "Senior" {
		return 4
	} else if profil == "Golang backend - developer" && position == "Team Lead" {
		return 5
	} else if profil == "Java backend - developer" && position == "Junior" {
		return 6
	} else if profil == "Java backend - developer" && position == "Intern" {
		return 7
	} else if profil == "Java backend - developer" && position == "Middle" {
		return 8
	} else if profil == "Java backend - developer" && position == "Senior" {
		return 9
	} else if profil == "Java backend - developer" && position == "Team Lead" {
		return 10
	} else if profil == "Frontend - developer" && position == "Junior" {
		return 11
	} else if profil == "Frontend - developer" && position == "Intern" {
		return 12
	} else if profil == "Frontend - developer" && position == "Middle" {
		return 13
	} else if profil == "Frontend - developer" && position == "Senior" {
		return 14
	} else if profil == "Frontend - developer" && position == "Team Lead" {
		return 15
	} else if profil == "Data Science - specialist" && position == "Junior" {
		return 16
	} else if profil == "Data Science - specialist" && position == "Intern" {
		return 17
	} else if profil == "Data Science - specialist" && position == "Middle" {
		return 18
	} else if profil == "Data Science - specialist" && position == "Senior" {
		return 19
	} else if profil == "Data Science - specialist" && position == "Team Lead" {
		return 20
	} else {
		return 3
	}
}
func isValidPhoneNumber(phoneNumber string) bool {
	// Регулярное выражение для проверки номера телефона
	// Поддерживаются форматы: 89387650971, +79387650971
	// Разрешены тире и пробелы между цифрами
	result := regexp.MustCompile(`^(8|\+7)[\s-]?\d{3}[\s-]?\d{3}[\s-]?\d{2}[\s-]?\d{2}$`)

	// Проверка соответствия номера телефона регулярному выражению
	return result.MatchString(phoneNumber)
}

type Config struct {
	TelegramBotToken string
	DbMysqlParams    string
}

func DecodeConfig(fileName string) (Config, error) {
	file, _ := os.Open(fileName)
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	return configuration, err
}
