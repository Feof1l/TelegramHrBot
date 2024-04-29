package telegramTest

import (
	"testing"
	"time"

	"github.com/Feof1l/TelegramHrBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestHandleUpdates(t *testing.T) {
	// Подготовка данных для теста
	mockUpdatesChan := make(chan tgbotapi.Update)
	bot := &telegram.Bot{
		//bot: nil, // Не требуется для этого теста
	}

	// Запуск обработки обновлений в отдельной горутине
	go bot.HandleUpdates(mockUpdatesChan)

	// Отправка пустого обновления для проверки обработки
	mockUpdatesChan <- tgbotapi.Update{}

	// Временная задержка для ожидания завершения обработки в горутине
	time.Sleep(100 * time.Millisecond)

	// Дополнительные проверки, если необходимо
	// Например, проверка состояния объекта bot после обработки обновления
}

// Дополнительные проверки, если необходимо
// Например, можно проверить, что бот был заблокирован

// Временная задержка для ожидания завершения обработки в горутине
// Обычно такие проверки делаются асинхронно, поэтому здесь используется временная задержка
// В реальном тесте это может быть сделано с помощью канала ожидания или другого механизма синхронизации
// В этом примере просто добавлена задержка для демонстрации

func TestHandleUpdates_Message_Start(t *testing.T) {
	// Подготовка данных для теста
	mockUpdatesChan := make(chan tgbotapi.Update)
	bot := &telegram.Bot{
		//bot: nil, // Не требуется для этого теста
	}

	// Запуск обработки обновлений в отдельной горутине
	go bot.HandleUpdates(mockUpdatesChan)

	// Отправка обновления с командой "start"
	mockUpdatesChan <- tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			Text: "/start",
		},
	}

	// Дополнительные проверки, если необходимо
	// Например, можно проверить, что было отправлено сообщение с клавиатурой

	// Временная задержка для ожидания завершения обработки в горутине
	// Обычно такие проверки делаются асинхронно, поэтому здесь используется временная задержка
	// В реальном тесте это может быть сделано с помощью канала ожидания или другого механизма синхронизации
	// В этом примере просто добавлена задержка для демонстрации
	time.Sleep(100 * time.Millisecond)
}
func TestHandleUpdates_Message_UnblockedUser(t *testing.T) {
	// Подготовка данных для теста
	mockUpdatesChan := make(chan tgbotapi.Update)
	bot := &telegram.Bot{
		//bot: nil, // Не требуется для этого теста
	}

	// Запуск обработки обновлений в отдельной горутине
	go bot.HandleUpdates(mockUpdatesChan)

	// Отправка обновления с текстовым сообщением от разблокированного пользователя
	mockUpdatesChan <- tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			Text: "Test message",
		},
	}

	// Дополнительные проверки, если необходимо
	// Например, можно проверить, что было отправлено сообщение

	// Временная задержка для ожидания завершения обработки в горутине
	time.Sleep(100 * time.Millisecond)
}
