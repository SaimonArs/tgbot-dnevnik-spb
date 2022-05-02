package dnevnik2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"main.go/lib/e"
	"main.go/lib/utilslice"
)

const (
    getListPeriod="group/group/get-list-period"
    getMarksTable="journal/estimate/table"
    getTimeTable="journal/lesson/list-by-education"
    getListStudents="journal/person/related-child-list"
    postLogin="user/auth/login"
)

var ErrAuth = errors.New("incorrect data")

type Client struct {
    host string
    basePath string
    client http.Client
}

func New(host, basePath string) *Client {
    return &Client{
        host: host,
        basePath: basePath,
        client: http.Client{},
    }
}

func (c *Client) TimeTable(token, pDate string, pEducations int) (TimeTable, error){
   const errMs = "err TimeTable"
    qu := url.Values{}                             
    qu.Add("p_datetime_to", fmt.Sprintf("%s 23:59:59", pDate))                   
    qu.Add("p_datetime_from", fmt.Sprintf("%s 00:00:00", pDate))       
    qu.Add("p_educations[]", strconv.Itoa(pEducations))  

    data, err := c.doReqeuest(getTimeTable, token, qu)
    if err != nil {
        return nil, e.Wrap(errMs, err) 
    }

    a := make(TimeTable)
    
    var res RestResponse
    if err := json.Unmarshal(data, &res); err != nil {
        return nil, e.Wrap(errMs, err)
    }

    for _, val := range res.Dat.Items {
        a[val.Number] = Lesson{
        Subject: val.Subject,
        Content: val.Content,
        HomeWork: val.Tasks,
        }
    }
    return a, nil
}

func (c *Client) MarksTable(token, pDateFrom, pDateTo string, pEducations, pLimit int) (SchoolMarks, error){
    const errMs = "err MarksTable"
    qu := url.Values{}
    qu.Add("p_date_to", pDateTo)
    qu.Add("p_date_from", pDateFrom)
    qu.Add("p_educations[]", strconv.Itoa(pEducations))
    qu.Add("p_limit", strconv.Itoa(pLimit))

    data, err := c.doReqeuest(getMarksTable, token, qu)
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }

    a := make(SchoolMarks)

    var res RestResponse
    if err := json.Unmarshal(data, &res); err != nil {
        return nil, fmt.Errorf("%s %s", errMs, err)
    }

    packingMarks(a, res)

    for i := 2; i <= res.Dat.Tpages; i++ {
        qu.Set("p_page", strconv.Itoa(i))
        data, err = c.doReqeuest(getMarksTable, token, qu)
        if err != nil {
            return nil, e.Wrap(errMs, err)
        }    
        
        if err := json.Unmarshal(data, &res); err != nil {
            return nil, e.Wrap(errMs, err)

        }

        packingMarks(a, res)
        
    }

    reverseMarks(a)

    return a, nil
}

func (c *Client) ListStudents(token string) (Student, error) { // !!!! Работает только после авторизации
    const errMs = "err get list students"
    qu := url.Values{}
    qu.Add("p_page", "1")

    data, err := c.doReqeuest(getListStudents, token, qu)
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }

    var res RestResponse

    if err = json.Unmarshal(data, &res); err != nil {
        return nil, e.Wrap(errMs, err)
    }

    a := make(Student)
    for _, val := range res.Dat.Items {       
        a[fmt.Sprintf("%s %s %s", val.Sname, val.Fname, val.Mname)] = val.Edu
    }
    return a, nil
}

func (c *Client) ListPeriod(token string, groupId int) ([]Period, error){
    const errMs = "err get list period"
    qu := url.Values{}
    qu.Add("p_page", "1")
    qu.Add("p_group_ids[]", strconv.Itoa(groupId))

    data, err := c.doReqeuest(getListPeriod, token, qu)
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }

    var res RestResponse
    if err = json.Unmarshal(data, &res); err != nil {
        return nil, e.Wrap(errMs, err)
    }

    a := []Period{}
    for _, val := range res.Dat.Items {
        a = append(a,
                Period{
                    Name: val.Name,
                    PdFrom: val.PdFrom,
                    PdTo: val.PdTo,
                })
        }
    return a, nil
}

func (c *Client) Authorization(password, email string) (string, error) {
    const errMs = "err get token: "
    u := url.URL{
        Scheme: "https",
        Host: c.host,
        Path: path.Join(c.basePath, postLogin),
    }

    message := map[string]interface{}{
        "activation_code":  nil,
        "login":            email,
        "password":         password, 
        "type":             "email",
        "_isEmpty":         false,
    }

    bRepre, err := json.Marshal(message)
    if err != nil {
        return "", e.Wrap(errMs, err)
    }

    req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(bRepre))
    if err != nil {
        return "", e.Wrap(errMs, err)
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Accept", "text/plain")
    req.Header.Add("Accept", "*/*")
    req.Header.Add("Access-Control-Allow-Origin", "*")
    req.Header.Add("Content-Type", "text/plain")
    req.Header.Add("User-Agent", "Go/application")

    
    resp, err := c.client.Do(req)
    if err != nil {
        return "", e.Wrap(errMs, err)
    }
    defer resp.Body.Close()    

    if resp.StatusCode / 100 == 2 {
        return resp.Header.Values("X-JWT-Token")[0], nil
    } else {
        return "", ErrAuth 
    }
}

func (c *Client) doReqeuest(method, token string, query url.Values) ([]byte, error) {
    const errMs = "err do request"
    u := url.URL{
        Scheme: "https",
        Host: c.host,
        Path: path.Join(c.basePath, method),
    }

    cookie := &http.Cookie{
        Name: "X-JWT-Token",
        Value: token,
    }

    req, err := http.NewRequest(http.MethodGet, u.String(), nil)
    req.URL.RawQuery = query.Encode()
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }

    req.Header.Add("User-Agent", "Go/application")

    req.AddCookie(cookie)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)

    if err != nil {
        return nil, e.Wrap(errMs, err)
    }
    return body, nil
}

func reverseMarks(a SchoolMarks) {
    for key, val := range a {
        a[key] = utilslice.ReversInt(val) 
    }
}

func packingMarks(a SchoolMarks, res RestResponse) {
    items := res.Dat.Items
    for _, val := range items {
        n, err := strconv.Atoi(val.Mark) 
        if err != nil {
            continue
        }

        _, ok := a[val.Subject]
        if !ok {
            a[val.Subject] = []int{n}
        } else {
            a[val.Subject] = append(a[val.Subject], n)
        }

    }
}
