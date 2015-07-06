package main

import (
    "fmt"
    "log"
    "os"
    "math/rand"
)

var logger = log.New(os.Stdout, "", log.Ltime | log.Lshortfile)

func main() {
    fmt.Println()
    LoginInstagram()
    posts := GetPosts("beyonce")

    for true {
        post_seed := rand.Intn(len(posts))
        var tmp []string
        for i, postcode := range posts {

            users := GetCommentators(postcode)
            user_seed := rand.Intn(len(users))
            for j, user := range users {

                if i==post_seed && j==user_seed {
                    tmp = GetPosts(user)
                } else {
                    GetPosts(user)
                }
            }
        }
        posts = tmp
    }
}
