package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"main.go/database"
	"main.go/lib/e"
)

type Database struct {
    sql *sql.DB
}

func New(DBfiles string) (*Database, error) {
    sqlDB, err := sql.Open("sqlite3", DBfiles)
    if err != nil {
        return nil, e.Wrap("err Open db", err)
    }

    if _, err = sqlDB.Exec(schemeSQL); err != nil {
        return nil, e.Wrap("err Exist tables", err)
    }

    db := Database{sql : sqlDB}
    return &db, nil
}

func (db Database) UserData(id int) (string, int, int, error) {
    stmt := db.sql.QueryRow("select * from users where id = $1", id)
    var token string
    var pedu, pgroup int
    err := stmt.Scan(&id, &pedu, &pgroup, &token)
    if err != nil {
        return "", 0, 0,  e.Wrap("err get user data", err)
    }
    return token, pedu, pgroup, nil
}

func (db Database) StuList(id int) ([]database.InfoEdu, error) {
    const errMs = "err students exist"
    stmt, err := db.sql.Query("select * from students where id = $1", id)
    if err != nil {
        return nil, e.Wrap(errMs, err)
    }
    defer stmt.Close()
    LS := []database.InfoEdu{}

    for stmt.Next(){
        s := database.InfoEdu{}
        err := stmt.Scan(&s.ID, &s.PEdu, &s.GroupID, &s.Info)
        if err != nil {
            continue
        }
        LS = append(LS, s)
    }
    return LS, nil
}

func (db Database) UpdateUser(u *database.User) error {
    stmt, err := db.sql.Prepare(updateUser)
    if err != nil {
        return e.Wrap("err Update token", err)
    }
    defer stmt.Close()
    stmt.Exec(u.Token, u.PEdu, u.GroupID, u.ID)
    
    return nil
}

func (db Database) CreateUser(u *database.User) error {
    stmt, err := db.sql.Prepare(insertUser)
    if err != nil {
        return e.Wrap("err insert user", err)
    }
    defer stmt.Close()
    stmt.Exec(u.ID, u.PEdu, u.GroupID, u.Token)
    return nil
}

func (db Database) DeleteUser(id int) error {
    stmt, err := db.sql.Prepare(deletUser)
    if err != nil {
        return e.Wrap("err delete user", err)
    }
    defer stmt.Close()
    stmt.Exec(id)
    return nil
}

func (db Database) ExistsUser(id int) (bool, error) {
    stmt := db.sql.QueryRow("select exists(select id from users where id = $1)", id) 
    var b bool  
    err := stmt.Scan(&b)
    if err != nil {
        return false, e.Wrap("err get state users", err)
    }
    return b, nil 
}

func (db Database) CreateStud(s *database.InfoEdu) error {
    stmt, err := db.sql.Prepare(insertStud)
    if err != nil {
        return e.Wrap("err insert students", err)
    }
    defer stmt.Close()
    stmt.Exec(s.ID, s.PEdu, s.GroupID, s.Info)
    return nil
}

func (db Database) DeleteStud(id int) error {
    stmt, err := db.sql.Prepare(deletStu)
    if err != nil {
        return e.Wrap("err delete students", err)
    }
    defer stmt.Close()
    stmt.Exec(id)
    return nil
}

func (db Database) ExistsStud(id int) (bool, error) {
    stmt := db.sql.QueryRow("select exists(select id from students where id = $1)", id) 
    var b bool  
    err := stmt.Scan(&b)
    if err != nil {
        return false, e.Wrap("err get state users", err)
    }
    return b, nil 
}
