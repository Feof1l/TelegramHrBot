package telegramTest

import (
	"log"
	"os"
	"testing"

	"github.com/Feof1l/TelegramHrBot/pkg/models/mysql"
	"github.com/Feof1l/TelegramHrBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestNewBot(t *testing.T) {
	// Подготовка данных для теста
	mockBotAPI := &tgbotapi.BotAPI{}
	mockErrorLog := log.New(os.Stdout, "error: ", log.LstdFlags)
	mockInfoLog := log.New(os.Stdout, "info: ", log.LstdFlags)
	mockCandidatesModel := &mysql.CandidatModel{} // замените на вашу реализацию, если требуется

	// Вызов конструктора для создания нового экземпляра бота
	bot := telegram.NewBot(mockBotAPI, mockErrorLog, mockInfoLog, mockCandidatesModel)

	// Проверка, что возвращенный экземпляр бота не равен nil
	if bot == nil {
		t.Error("Ожидался не nil экземпляр бота, но получен nil")
	}

	// Дополнительные проверки, если необходимо
}
