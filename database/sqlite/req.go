package sqlite

const schemeSQL = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY, 
    pedu INTEGER NOT NULL,
    pgroup INTEGER NOT NULL,
    token TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS students (
    id INTEGER NOT NULL,
    pedu INTEGER NOT NULL,
    pgroup INTEGER NOT NULL,
    info INFO NOT NULL
)
`

const insertStud = `
INSERT INTO students (id, pedu, pgroup, info) 
values (?, ?, ?, ?)
`

const insertUser = `
INSERT INTO users (id, pedu, pgroup, token) 
values (?, ?, ?, ?)
`

const deletUser = `
DELETE FROM users
where id = ?
`

const deletStu = `
DELETE FROM students 
where id = ?
`

const updateUser = `
UPDATE users 
set token = ?, pedu = ?, pgroup = ?
where id = ?
`


