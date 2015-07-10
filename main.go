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

    Login()

    user_id, posts := GetPosts("natgeo")

    for true {
        var seed_posts []string
        for _, post := range posts {
            users := GetCommentators(post)
            for _, user := range users {
                user_id, seed_posts = GetPosts(user)
                if len(posts)>5 && !Check(user) && user_id!="" {
                    Follow(user_id, user)
                    break;
                }
            }
        }
        posts = seed_posts
    }
}
