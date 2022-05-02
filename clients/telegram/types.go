package  telegram

type Update struct {
    UpdateId int       `json:"update_id"`
    Message *Message    `json:"message"`
    Callback *Callback  `json:"callback_query"`
}

type Message struct {
    User User           `json:"from"`
    MsgId int           `json:"message_id"` 
    Chat Chat           `json:"chat"`
    Text string         `json:"text"`
}

type Callback struct {
    Message Message     `json:"message"`
    Data string         `json:"data"`
}

type User struct {
    Username string `json:"username"`
}

type Chat struct {
    ChatId int          `json:"id"`
}

type RestResponse struct {
    Ok bool             `json:"ok"`
    Result []Update     `json:"result"`
}
 
type BotMessage struct {
    ChatId int          `json:"chat_id"`
    Text string         `json:"text"`
}

type InlineKeyboardButton struct {
    Text string `json:"text"`
    Callback_data string `json:"callback_data"`
}
type InlineKeyboardMarkup struct {
    Inline_keyboard [][]InlineKeyboardButton  `json:"inline_keyboard"`
}
