package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Feof1l/TelegramHrBot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Config struct {
	TelegramBotToken string
}

func main() {

	// Определение нового флага из командной строки для настройки MySQL подключения.
	dsn := flag.String("dsn", "admin:dfyz10012003dfyz@/candidates?parseTime=true", "Название MySQL источника данных")
	// извлекаем флаг из командной строки
	flag.Parse()
	// создаем новый логер для вывода информационных сообщенйи в поток stdout c припиской info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//аналогично для логов с ошибками, такеж включим вывод фйла и номера  строки, где произошла ошибка
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Мы также откладываем вызов db.Close(), чтобы пул соединений был закрыт
	// до выхода из функции main().
	// Подробнее про defer: https://golangs.org/errors#defer
	defer db.Close()

	// декодируем файл json, в котором хранится конфиг - токен бота

	configuration, err := decodeConfig("config.json")
	if err != nil {
		errorLog.Println(err)
	}
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		errorLog.Println(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, errorLog, infoLog)
	func() {
		if err := telegramBot.Start(); err != nil && !telegramBot.IsBlockedUser() { // потом переделать
			errorLog.Println(err)
		}
	}()

	////////////////////////

}
func decodeConfig(fileName string) (Config, error) {
	file, _ := os.Open(fileName)
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	return configuration, err
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
