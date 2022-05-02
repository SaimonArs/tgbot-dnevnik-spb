package telegram

const msgHelp =
`*Для авторизации введи:* /reg почта пароль _(Данные от дневника!!!)_
*Для удаления данных из бота:* /logout
*Для выбора обучающегося:* /switch
*Для просмотра оценок:* /marks 
*Для просмотра расписания:* /timetable Дата _(Формат даты дд.мм.гггг)_
`


const msgHello = 
`*Telegram дневник Петербургского образования*,
_неофициальный бот для сайта https://dnevnik2.petersburgedu.ru_` + "\n\n"+ msgHelp
const (
    msgUnknownCommand = "Неизвестная команда"
)
