package database 

import (
	"errors"
)

type Database interface {
    CreateUser(u *User) error
    CreateStud(s *InfoEdu) error
    UpdateUser(u *User) error
    DeleteUser(id int) error
    DeleteStud(id int) error
    UserData(id int) (string, int, int, error)
    StuList(id int) ([]InfoEdu, error)
    ExistsUser(id int) (bool, error)
    ExistsStud(id int) (bool, error)
}

var ErrNoLogin = errors.New("no loged user")

type User struct {
    ID int
    Token string
    PEdu int
    GroupID int
}

type InfoEdu struct {
    ID int
    PEdu int
    GroupID int
    Info string
}
