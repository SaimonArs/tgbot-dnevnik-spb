package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"main.go/database"
	"main.go/lib/e"
)

func (p *Processor) doCallback(text string, chatId, msgId int) error {
    var ftext []string
    text = strings.TrimSpace(text)
    for _, val := range strings.Split(text, " ") {
        if val != "" {
            ftext = append(ftext, val)
        }
    }
    p.tg.DeleteMessage(chatId, msgId)
    
    log.Printf("got new callback %s ", ftext[0])
    switch ftext[0] {
    case "switch":
        return p.switchPedu(chatId, ftext)
    case "marks":
        return p.marks(chatId, ftext)
    default:
        return nil
    }
}

func (p *Processor) marks(chatId int, text []string) error {
    const errMs = "err marks"
    if len(text) != 3 {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return nil
    }
    
    t, edu, _, err := p.storage.UserData(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)
    }
    a, err := p.dnevnik.MarksTable(t, text[1], text[2], edu, 100)
    answer := []string{}
    if err != nil {
        return e.Wrap(errMs, err)
    }
    var am float32 = 0 
    c := []string{}
    for key, val := range a {
        for _, val2 := range val {
            am += float32(val2)
            c = append(c, strconv.Itoa(val2))
        }
        answer = append(answer, fmt.Sprintf("*%s*:\n _Ср. Арифм. %1.2f_ \n %s", key, am / float32(len(val)), strings.Join(c, " ")))
        am = 0
        c = []string{}
    }
    p.tg.SendMessage(chatId, strings.Join(answer, "\n\n"), nil)
    return nil
}

func (p *Processor) switchPedu(chatId int, text []string) error {
    const errMs = "err switchEdu"

    a, err := p.storage.StuList(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)
    } 
    numb, err := strconv.Atoi(text[1])
    if err != nil {
        p.tg.SendMessage(chatId, "Внутренняя ошибка", nil)
        return e.Wrap(errMs, err)
    }
    if len(a) - 1 < numb || numb < 0 {
        p.tg.SendMessage(chatId, "Этот номер отсутствует", nil)
        return nil
    }

    t, _, _, err := p.storage.UserData(chatId)
    if err != nil {
        p.tg.SendMessage(chatId, "Ошибка обновления", nil)
        return e.Wrap(errMs, err)
    }
    err = p.storage.UpdateUser(&database.User{ID: chatId, PEdu: a[numb].PEdu, GroupID: a[numb].GroupID, Token: t})
    if err != nil {
        p.tg.SendMessage(chatId, "Ошибка обновления", nil)
        return e.Wrap(errMs, err)
    }

    p.tg.SendMessage(chatId, "Информация обновлена", nil)
    return nil 
}
