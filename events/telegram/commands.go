package telegram

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"main.go/clients/dnevnik2"
	"main.go/clients/telegram"
	"main.go/database"
	"main.go/lib/e"
)

const (
    HelpCmd = "/help"
    StartCmd = "/start"
    LoginCmd = "/reg"
    DeleteCmd = "/logout"
    SwitchCmd = "/switch"
    MarksCmd = "/marks"
    TableCmd = "/timetable"
)

func (p *Processor) doCmd(text string, chatId, msgId int, username string) error {
    var ftext []string
    text = strings.TrimSpace(text)
    for _, val := range strings.Split(text, " ") {
        if val != "" {
            ftext = append(ftext, val)
        }
    }
    
    log.Printf("got new command %s from %s", ftext[0], username)
    
    switch ftext[0] { 
    case LoginCmd:
        return p.login(chatId, msgId, ftext)
    case DeleteCmd:
        return p.deleting(chatId)
    case SwitchCmd:
        return p.showStudents(chatId)
    case MarksCmd:
        return p.marksPeriod(chatId)
    case TableCmd:
        return p.table(chatId, ftext)
    case HelpCmd:
        return p.sendHelp(chatId)
    case StartCmd:
        return p.sendHello(chatId)
    default:
        return p.tg.SendMessage(chatId, msgUnknownCommand, nil) 
    }
}

func (p *Processor) marksPeriod(chatId int) error {
    const errMs = "err marks period"

    if ok, _ := p.storage.ExistsUser(chatId); !ok {
        p.tg.SendMessage(chatId, "Сначала авторизуйтесь!", nil)
        return nil 
    }

    t, _, group, err := p.storage.UserData(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)
    }

    answer, err := p.dnevnik.ListPeriod(t, group)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)
    }

    keyb := telegram.InlineKeyboardMarkup{Inline_keyboard: [][]telegram.InlineKeyboardButton{}}

    for _, val := range answer {
        keyb.Inline_keyboard = append(keyb.Inline_keyboard, []telegram.InlineKeyboardButton{{ Text: val.Name, Callback_data: "marks " + val.PdFrom + " " + val.PdTo}})    
    }

    keyb.Inline_keyboard = append(keyb.Inline_keyboard, []telegram.InlineKeyboardButton{{Text: "Отмена", Callback_data: "exit"}})

    p.tg.SendMessage(chatId, "Выберете период:", keyb)
    return nil
}

func (p *Processor) table(chatId int, text []string) error {
    const errMs = "err time table"
    if len(text) != 2 {
        p.tg.SendMessage(chatId, `Неправильный формат ввода, /timetable дд.мм.гггг`, nil)
        return nil
    }

    if ok, _ := p.storage.ExistsUser(chatId); !ok {
        p.tg.SendMessage(chatId, "Сначала авторизуйтесь!", nil)
        return nil 
    }

    t, edu, _, err := p.storage.UserData(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)
    }
    ai := []int{}
    answer := []string{}
    a, err := p.dnevnik.TimeTable(t, text[1], edu)
    if err != nil {
        p.tg.SendMessage(chatId, `Неправильный формат ввода, /timetable дд.мм.гггг`, nil)
        return e.Wrap(errMs, err)
    }

    for idx := range a {
        ai = append(ai, idx)
    }
    sort.Ints(ai)

    hw := []string{}
    for _, idx := range ai {
        for _, val := range a[idx].HomeWork {
            hw = append(hw, val.Tname)
        }
        answer = append(answer, fmt.Sprintf("*%d: %s*\n%s\n*Домашнее Задание:* \n _%s_", idx, a[idx].Subject, a[idx].Content, strings.Join(hw, "\n")))
        hw = []string{}
    }
    str := strings.Join(answer, "\n\n")
    if (str == "") {
        p.tg.SendMessage(chatId, "*Домашнее задание на этот день отсутствует*", nil)
    } else {
        p.tg.SendMessage(chatId, str, nil)
    }
    return nil
}



func (p *Processor) login(chatId, msgId int, text []string) error {
    const errMs = "err login" 

    p.tg.DeleteMessage(chatId, msgId)

    if len(text) != 3 {
        p.tg.SendMessage(chatId, `Неправильный формат ввода, /reg почта пароль`, nil)
        return nil
    }
    
    b, err := p.storage.ExistsStud(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)    
    }
    if b {
        p.tg.SendMessage(chatId, "Вы уже зарегистрированы", nil)
        return nil
    }

    token, err := p.dnevnik.Authorization(text[2], text[1])
    if err != nil {
        p.tg.SendMessage(chatId, "Ошибка при авторизации", nil)
        return e.Wrap(errMs, err) 
    }

    
    var list dnevnik2.Student
    list, err = p.dnevnik.ListStudents(token)
    if err != nil {
        return e.Wrap(errMs, err) 
    }
    
    pe := []int{}
    pg := []int{}
    
    for key, val := range list {
        for _, val2 := range val {
            s := database.InfoEdu{
                ID: chatId,
                PEdu: val2.EduID,
                GroupID: val2.GroupID,
                Info: val2.InstName + " : " + key,
            }
            
            pe = append(pe, val2.EduID)
            pg = append(pg, val2.GroupID)
            err = p.storage.CreateStud(&s)
            if err != nil {
                return e.Wrap(errMs, err) 
            }
        }
    }
    if len(pe) == 0 || len(pg) == 0 {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)
    }
    u := database.User{
        ID: chatId,
        PEdu: pe[0],
        GroupID: pg[0],
        Token: token,
    }
    p.storage.CreateUser(&u)

    p.tg.SendMessage(chatId, "Регистрация прошла успешно!", nil)

    return nil
}

func (p *Processor) showStudents(chatId int) error {
    const errMs = "err show Students" 
    if ok, _ := p.storage.ExistsUser(chatId); !ok {
        p.tg.SendMessage(chatId, "Сначала авторизуйтесь!", nil)
        return nil 
    }
    if ok, _ := p.storage.ExistsStud(chatId); !ok {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return nil 
    }
    a, err := p.storage.StuList(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err) 
    }
    keyb := telegram.InlineKeyboardMarkup{ Inline_keyboard: [][]telegram.InlineKeyboardButton{}}
    for idx, val := range a {
        Aidx := strconv.Itoa(idx)
        keyb.Inline_keyboard = append(keyb.Inline_keyboard, []telegram.InlineKeyboardButton{{Text: val.Info, Callback_data: "switch " + Aidx}})
    }
    
    keyb.Inline_keyboard = append(keyb.Inline_keyboard, []telegram.InlineKeyboardButton{{Text: "Отмена", Callback_data: "exit"}})
    p.tg.SendMessage(chatId, "Выберите обучающегося:", keyb)
    return nil
}


func (p *Processor) deleting(chatId int) error {
    const errMs = "err deleting"
    err := p.storage.DeleteUser(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err) 
    }
    
    err = p.storage.DeleteStud(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err) 
    }

    p.tg.SendMessage(chatId, "Удаление данных прошло успешно", nil)
    return nil  
}

func (p *Processor) sendHelp(chatId int) error {
    return p.tg.SendMessage(chatId, msgHelp, nil)
}

func (p *Processor) sendHello(chatId int) error {
    return p.tg.SendMessage(chatId, msgHello, nil)
}
