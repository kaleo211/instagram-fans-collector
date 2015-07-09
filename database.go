package main

import (
    _ "github.com/go-sql-driver/mysql"
)

func CheckErr(err error) {
    if err!=nil {
        logger.Println(err)
    }
}

func Next() (name string) {
    rows, err := db.Query("SELECT username FROM tofollow LIMIT 1")
    defer rows.Close()
    CheckErr(err)

    if rows.Next() {
        rows.Scan(&name)
        stmt, _ := db.Prepare("DELETE FROM tofollow WHERE username=? LIMIT 1")
        defer stmt.Close()
        _, err = stmt.Exec(name)
        CheckErr(err)
        stmt.Close()
    }
    return
}

func SaveToFollow(username string) {
    rows, err := db.Query("SELECT * FROM tofollow WHERE username=?", username)
    defer rows.Close()
    CheckErr(err)
    if rows.Next() {
        return
    }

    stmt, _ := db.Prepare("INSERT tofollow SET username=?")
    defer stmt.Close()
    _, err = stmt.Exec(username)
    CheckErr(err)
}


func SaveFollowed(username string) {
    stmt, _ := db.Prepare("INSERT users SET username=?")
    defer stmt.Close()
    _, err := stmt.Exec(username)
    CheckErr(err)
}

func Check(username string) bool {
    rows, err := db.Query("SELECT * FROM users WHERE username = ? ", username)
    defer rows.Close()
    CheckErr(err)
    if rows.Next() {
        return true
    }
    return false
}
