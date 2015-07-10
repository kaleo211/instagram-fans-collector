package main

import (
    "bytes"
    "net/http"
    "net/url"
    "golang.org/x/net/html"
    "strings"
    "regexp"
    "time"
    "io/ioutil"
)


var client = &http.Client{}
var transport = &http.Transport{}
var cookies = make([]*http.Cookie, 20)
var cookies_size = 0

const account = "dirty.lily"
const password = "zxcvbnm123"
const host = "https://instagram.com"


func UpdateCookies(cokies []*http.Cookie) {
    for _, c := range cokies {
        updated := false
        for i:=0; i<cookies_size; i+=1 {
            if cookies[i].Name==c.Name {
                cookies[i].Value = c.Value
                updated = true
            }
        }
        if !updated {
            cookies[cookies_size] = &http.Cookie{Name: c.Name, Value: c.Value}
            cookies_size += 1
        }
    }
}

func Login() {

    instagram := "https://instagram.com/accounts/login/"
    request, _ := http.NewRequest("GET", instagram, nil)
    response, _ := client.Do(request)
    UpdateCookies(response.Cookies())

    // set post form data
    data := url.Values{}
    data.Set("username", account)
    data.Set("password", password)

    instagram_login := "https://instagram.com/accounts/login/ajax/"
    req, _ := http.NewRequest("POST", instagram_login, bytes.NewBufferString(data.Encode()))

    for i:=0; i<cookies_size; i+=1 {
        c := cookies[i]
        if c.Name=="csrftoken" {
            req.AddCookie(c)
            req.Header.Set("X-CSRFToken", c.Value);
        }
        if c.Name=="mid" {
            req.AddCookie(c)
        }
    }

    req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Accept-Language", "en-US,en;q=0.5")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Content-Length", "44")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
    req.Header.Set("Host", "instagram.com")
    req.Header.Set("Pragma", "no-cache")
    req.Header.Set("Referer", "https://instagram.com/accounts/login/")
    req.Header.Set("User-Agent", "  Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:36.0) Gecko/20100101 Firefox/36.0")
    req.Header.Set("X-Instagram-AJAX", "1")
    req.Header.Set("X-Requested-With", "XMLHttpRequest")

    resp, _ := client.Do(req)
    defer resp.Body.Close()
    UpdateCookies(resp.Cookies())

    contents, _ := ioutil.ReadAll(resp.Body)
    logger.Println("%s", string(contents))

    if resp.StatusCode==200 {
        logger.Println("login into Instagram successfully.")
    } else {
        logger.Fatalln("failed login into Instagram.")
    }
}


func Logout() {

    logout := "https://instagram.com/accounts/logout/"
    req, _ := http.NewRequest("POST", logout, bytes.NewBufferString(url.Values{}.Encode()))

    for i:=0; i<cookies_size; i+=1 {
        c := cookies[i]
        if c.Name=="csrftoken" || c.Name=="mid" || c.Name=="sessionid" || c.Name=="ds_user_id" {
            req.AddCookie(c)
        }
    }

    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Accept-Language", "en-US,en;q=0.5")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Host", "instagram.com")
    req.Header.Set("Referer", "https://instagram.com/dirty.lily/")
    req.Header.Set("User-Agent", "  Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:36.0) Gecko/20100101 Firefox/36.0")

    resp, _ := client.Do(req)
    defer resp.Body.Close()
    UpdateCookies(resp.Cookies())

    logger.Println(resp.StatusCode)
}

func GetPosts(username string) (user_id string, posts []string) {

    posts_url := host +"/" + username
    req, _ := http.NewRequest("GET", posts_url, nil)
    for i:=0; i<cookies_size; i+=1 {
        c := cookies[i]
        req.AddCookie(c)
    }

    resp, _ := transport.RoundTrip(req)
    defer resp.Body.Close()

    if resp.StatusCode!=200 {
        logger.Println("Failed to retrieve posts for", username)
        return
    }
    UpdateCookies(resp.Cookies())
    doc, _ := html.Parse(resp.Body)

    var data string
    var f func(*html.Node)
    // high level traverse html
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "script" {
            if child := n.FirstChild; child!=nil && strings.HasPrefix(child.Data, "window._sharedData = ") {
                data = strings.TrimPrefix(child.Data, "window._sharedData = ")
                data = strings.TrimSuffix(data, ";")
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)

    re_user_id := regexp.MustCompile("\"owner\"[^\"]*\"id\"[^\"]*\"([0-9]+)\"")
    match := re_user_id.FindStringSubmatch(data)
    if len(match) > 0 {
        user_id = match[1]
    }

    array := make([]string, 100)
    re_postcode := regexp.MustCompile("{\"code\":\"([\\w]{10})\"")
    matches := re_postcode.FindAllStringSubmatch(data, -1)
    for i, e := range matches {
        if i>=100 {break;}
        array[i] = e[1]
    }
    posts = array[0:len(matches)]
    logger.Println("Retrieved", len(posts), "posts for user:", username)

    return
}


func GetCommentators(postcode string) (users []string) {
    comments_url := host + "/p/" + postcode + "/"
    req, _ := http.NewRequest("GET", comments_url, nil)
    for i:=0; i<cookies_size; i+=1 {
        req.AddCookie(cookies[i])
    }

    resp, _ := transport.RoundTrip(req)
    defer resp.Body.Close()

    if resp.StatusCode!=200 {
        logger.Println("Failed to retrieve commentators for post:", postcode)
        users = make([]string, 0)
        return
    }

    UpdateCookies(resp.Cookies())
    doc, _ := html.Parse(resp.Body)

    var data string
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "script" {
            if child := n.FirstChild; child!=nil && strings.HasPrefix(child.Data, "window._sharedData = ") {
                data = strings.TrimPrefix(child.Data, "window._sharedData = ")
                data = strings.TrimSuffix(data, ";")
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)

    commentators := make([]string, 100)
    re := regexp.MustCompile("{\"username\":\"([\\w]+)\"")
    matches := re.FindAllStringSubmatch(data, -1)
    for i, e := range matches {
        if (i>=100) {break;}
        commentators[i] = e[1]
    }
    users = commentators[:len(matches)]
    logger.Println("Retrieved", len(users), "commentators for post:", postcode)

    return
}


func Follow(user_id string, username string) {

    if (Check(username)) {
        logger.Println(username, "had already been followed.")
        return
    }

    data := url.Values{}
    like_url := "https://instagram.com/web/friendships/" + user_id + "/follow/"
    req, _ := http.NewRequest("POST", like_url, bytes.NewBufferString(data.Encode()))

    for i:=0; i<cookies_size; i+=1 {
        c := cookies[i]
        if c.Name=="csrftoken" || c.Name=="mid" || c.Name=="sessionid" || c.Name=="ds_user_id" {
            req.AddCookie(c)
        }
        if c.Name=="csrftoken" {
            req.Header.Set("X-CSRFToken", c.Value)
        }
    }
    req.AddCookie(&http.Cookie{Name: "ig_pr", Value: "2"})
    req.AddCookie(&http.Cookie{Name: "ig_vw", Value: "1429"})
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Accept-Language", "en-US,en;q=0.5")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Content-Length", "0")
    req.Header.Set("Host", "instagram.com")
    req.Header.Set("Pragma", "no-cache")
    req.Header.Set("Referer", "https://instagram.com/"+username+"/")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:39.0) Gecko/20100101 Firefox/39.0")
    req.Header.Set("X-Instagram-AJAX", "1")
    req.Header.Set("X-Requested-With", "XMLHttpRequest")

    resp, _ := transport.RoundTrip(req)
    defer resp.Body.Close()
    if resp.StatusCode==200 {
        logger.Println("Followed", username+"("+user_id+") successfully.")
        SaveFollowed(username)
        return
    }

    logger.Println(resp.Status)
    logger.Println("Failed to follow", username+"("+user_id+")")
    if resp.StatusCode==403 {
        logger.Println("Sleep for 70 seconds...Then relogin")
        time.Sleep(70 * time.Second)
        // Logout()
        // Login()
    }
}
