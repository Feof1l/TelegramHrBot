package telegramTest

import (
	"testing"
	"time"

	"github.com/Feof1l/TelegramHrBot/pkg/telegram"
	"github.com/stretchr/testify/assert"
)

func TestClearChatHistoryPerformance(t *testing.T) {
	// Создайте экземпляр бота для тестирования
	bot := &telegram.Bot{}

	// Создайте фиктивные данные для чата и сообщений

	// Измерьте время выполнения метода clearChatHistory
	startTime := time.Now()
	err := bot.ClearChatHistory(653924346)
	duration := time.Since(startTime)

	// Проверьте, что метод завершился без ошибок и выполнение заняло разумное время
	assert.NoError(t, err, "clearChatHistory method returned an error")
	assert.True(t, duration < time.Second, "clearChatHistory method took too long: %v", duration)
}
