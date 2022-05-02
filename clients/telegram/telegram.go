package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"main.go/lib/e"
)

const(
    getUpdateMethod="getUpdates"
    sendMessageMethod="sendMessage"
    deleteMessageMethod="deleteMessage"

)

type Client struct {
    host string
    basePath string
    client http.Client
}

func New(host, token string) *Client {
    return &Client{
        host:   host,
        basePath: newBasePath(token),
        client: http.Client{},
    }
}

func (c *Client) SendMessage(chatId int, text string, inline interface{}) error {

    qu := url.Values{}
    qu.Add("chat_id", strconv.Itoa(chatId))
    qu.Add("parse_mode", "Markdown")
    qu.Add("text", text)

    if inline != nil {
        s, err := json.Marshal(inline)
        if err != nil {
            return e.Wrap("err send inline markup", err)
        }
        qu.Add("reply_markup", string(string(s)))
    }

    _, err := c.doReqeuest(sendMessageMethod, qu)
    if err != nil {
        return e.Wrap("err send message", err)
    }
    return nil
}

func (c *Client) DeleteMessage(chatId int, messageId int) error {
    qu := url.Values{}
    qu.Add("chat_id", strconv.Itoa(chatId))
    qu.Add("message_id", strconv.Itoa(messageId))

    _, err := c.doReqeuest(deleteMessageMethod, qu) 
    if err != nil {
        return e.Wrap("err edit message", err)
    }
    return nil
}

func (c *Client) Updates(offset int, limit int) ([]Update, error){ 
    const errMs = "err Updates"
    qu := url.Values{}
    qu.Add("offset", strconv.Itoa(offset))
    qu.Add("limit", strconv.Itoa(limit))

    data, err := c.doReqeuest(getUpdateMethod, qu)
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }

    var res RestResponse

    if err:= json.Unmarshal(data, &res); err != nil{
        return nil, e.Wrap(errMs, err)
    }

    return res.Result, nil

}

func (c *Client) doReqeuest(method string, query url.Values) ([]byte, error) {
    const errMs = "err do request"
    u := url.URL{
        Scheme: "https",
        Host: c.host,
        Path: path.Join(c.basePath, method),
    }

    req, err := http.NewRequest(http.MethodGet, u.String(), nil)
    req.URL.RawQuery = query.Encode()
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    return body, nil
}

func newBasePath(token string) string {
    return "bot" + token
}
