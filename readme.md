## Telegram Hr бот
Главная цель данного бота - оптимизация работы HR-ов.С помощью бота сократится время общения сотрудника (живого человека) с кандидатом на вакансию.Вместо этого кандидат будет вести диалог-переписку посредством соц-сети мессенджера(telegram) с ботом.Также,необходимые данные будут автоматически заносится в БД.
Данный бот позволит отсеивать несоответствующих  требованиям компании кандидатов, что сократит рабочее время и силы сотрудника, проводящего собеседования. Данный сотрудник будет вести беседу и проверять знания у кандидатов, прошедших переписку с ботом, то есть с теми, кого бот не отсеил.Таким образом, данный кандидат отвечает минимальным требованиям компании-нанимателя.


Версия Go : 1.20

### Конфиг
Токен для запуска бота находится в json - файле
Для запуска требуется получить токен и поместить в аналогичный файл


### Запуск приложения и работа с терминалом
Для запуска использовать команду:
```
$make run
```


### Логирование
Логирование осуществляется в два потока: информационные сообщения выводятся в поток stdout с префиксом INFO, сообщения об ошибках - в поток stderr, префикс - ERROR.
Можно перенаправить потоки из stdout и stderr в файлы на диске при запуске приложения из терминала следующим образом:
```
$ go run ./cmd/bot >>/tmp/info.log 2>>/tmp/error.log
```