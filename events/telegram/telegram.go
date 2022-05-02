package telegram

import (
	"errors"

	"main.go/clients/dnevnik2"
	"main.go/clients/telegram"
	"main.go/database"
	"main.go/events"
	"main.go/lib/e"
)

type Processor struct {
    tg *telegram.Client
    offset int
    storage database.Database
    dnevnik dnevnik2.Client
}

type Meta struct {
    ChatId int
    MsgId int
    Username string
}

var ( 
    ErrUnknownEventType = errors.New("unknown event type")
    ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage database.Database, dnevnik dnevnik2.Client) *Processor{
    return &Processor{
        tg: client,
        storage: storage,
        dnevnik: dnevnik,
    }
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
    updates, err := p.tg.Updates(p.offset, limit)
    if err != nil {
        return nil, e.Wrap("err get events", err)
    }

    if len(updates) == 0 {
        return nil, nil
    }

    res := make([]events.Event, 0, len(updates))

    for _, u := range updates {
        res = append(res, event(u))
    }

    p.offset = updates[len(updates)-1].UpdateId + 1

    return res, nil
}

func (p *Processor) Process(event events.Event) error {
    switch event.Type {
    case events.Message:
        return p.processMessage(event)
    case events.Inline:
        return p.processCallback(event)
    default:
        return e.Wrap("err process message", ErrUnknownEventType)
    }
}

func (p *Processor) processMessage(event events.Event) error {
    const errMs = "err process message"
    meta, err := meta(event)
    if err != nil {
        return e.Wrap(errMs, err)
    }

    if err := p.doCmd(event.Text, meta.ChatId, meta.MsgId, meta.Username); err != nil {
        return e.Wrap(errMs, err)
    }
    return nil
}

func (p *Processor) processCallback(event events.Event) error {
    const errMs = "err process callback"
    meta, err := meta(event)
    if err != nil {
        return e.Wrap(errMs, err)
    }
    if err := p.doCallback(event.Text, meta.ChatId, meta.MsgId); err != nil {
        return e.Wrap(errMs, err)
    }
    return nil

}

func meta(event events.Event) (Meta, error) {
    res, ok := event.Meta.(Meta)
    if !ok {
        return Meta{}, e.Wrap("err, get meta", ErrUnknownMetaType)
    }
    return res, nil
}

func event(upd telegram.Update) events.Event {
    udpType := fetchType(upd)
    res := events.Event{
        Type: udpType,
        Text: fetchText(upd),
    }
    if udpType == events.Message {
        res.Meta = Meta{
            ChatId: upd.Message.Chat.ChatId,
            MsgId: upd.Message.MsgId,
            Username: upd.Message.User.Username,
        }
    }
    if udpType == events.Inline {
        res.Meta = Meta{
            ChatId: upd.Callback.Message.Chat.ChatId,
            MsgId: upd.Callback.Message.MsgId,
            Username: upd.Callback.Message.User.Username,
        }
    }
    return res
}

func fetchText(upd telegram.Update) string{
   if upd.Message != nil && upd.Message.Text != ""{
       return   upd.Message.Text 
   }
   if upd.Callback != nil {
       return upd.Callback.Data
   }
   return ""
}

func fetchType(upd telegram.Update) events.Type {
    if upd.Message != nil && upd.Message.Text != ""{
        return events.Message
    }
    if upd.Callback != nil {
        return events.Inline
    } 
    return events.Unknown
}
