package dnevnik2

type SchoolMarks map[string] []int 

type TimeTable map[int] Lesson 

type Student map[string] []Educations
    
type Lesson struct {
    Subject string
    Content string
    HomeWork []Task 
}

type Task struct {
    Tname string        `json:"task_name"`
}

type Educations struct{
    EduID int           `json:"education_id"` 
    GroupID int         `json:"group_id"`
    InstName string     `json:"institution_name"`
}

type Period struct{
    Name string        
    PdFrom string       
    PdTo string         
}

type Items struct {
    Edu []Educations    `json:"educations"`
    Number int          `json:"number"`
    Subject string      `json:"subject_name"`
    Date string         `json:"date"`
    Mark string         `json:"estimate_value_name"`
    Content string      `json:"content_name"`
    Tasks []Task        `json:"tasks"`
    Fname string        `json:"firstname"`
    Sname string        `json:"surname"`
    Mname string        `json:"middlename"`
    Name string         `json:"name"`
    PdFrom string       `json:"date_from"`
    PdTo string         `json:"date_to"`
}

type data struct {
    Items []Items       `json:"items"` 
    Current int         `json:"current"`
    Tpages int          `json:"total_pages"`
}
 
type RestResponse struct {
    Dat data            `json:"data"`
}
