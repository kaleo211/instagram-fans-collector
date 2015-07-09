package main

import (
    "log"
    "os"
    "database/sql"
)

var logger = log.New(os.Stdout, "", log.Ltime | log.Lshortfile)
var db *sql.DB

func main() {

    db,_ = sql.Open("mysql", "kaleo211:iampassword@/instagram")

    LoginInstagram()

    SaveToFollow("natgeo")

    for true {
        user := Next()
        posts := GetPosts(user)
        for _, post := range posts {
            users := GetCommentators(post)
            for _, user = range users {
                SaveToFollow(user)
            }
        }
    }
}
