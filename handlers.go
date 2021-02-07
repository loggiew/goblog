// +build ignore

package main

import (
    "fmt"
    "time"
    "strconv"
    "database/sql"
    "net/http"
    "html/template"
    _ "github.com/go-sql-driver/mysql"
    "github.com/satori/go.uuid"
    "github.com/gorilla/mux"
)


type Article struct {
    Id int `json:"Id"`
    Title string `json:"Title"`
    Auth string `json:"Author"`
    Description sql.NullString `json:"Description"`
    Content string `json:"Content"`
}

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!<BR>")
    fmt.Fprintf(w, "Hi there, I love %s!<BR>", r.URL.Path[1:])
    fmt.Println("Endpoint Hit: homePage")
}


func test(w http.ResponseWriter, r *http.Request){
    blah := 1
    tmpl := template.Must(template.ParseFiles("template/test.html"))
    tmpl.Execute(w, blah)
}

func LoginHandler(w http.ResponseWriter, r *http.Request){

}

func AllArticlesGet(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    tmpl := template.Must(template.ParseFiles("template/articles.html"))
    stuff := ArticleLookup()
    //fmt.Printf("%#v\n", stuff)
    tmpl.Execute(w, stuff)
    //json.NewEncoder(w).Encode(Articles)
}

func CreatePost(w http.ResponseWriter, r *http.Request){
    tmpl := template.Must(template.ParseFiles("template/createpost.html"))
    tmpl.Execute(w, nil)
}

func AllArticlesPost(w http.ResponseWriter, r *http.Request){
    // Sessions
    response, err := checkSession(w, r)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        tmpl := template.Must(template.ParseFiles("template/unauthorized.html"))
        tmpl.Execute(w, nil)
        return
    }
    //content, err := strconv.Atoi(mux.Vars(r)["content"])
    //content := mux.Vars(r)["content"]
    var test sql.NullString
    test = sql.NullString{String:r.FormValue("description"), Valid: true}
    fmt.Println(test)
    content := r.FormValue("content")
    title := r.FormValue("title")
    description := sql.NullString{String:r.FormValue("description"), Valid: true}

    username := fmt.Sprintf("%s", response)
    var article Article
    article.Auth = username
    article.Content = content
    article.Title = title
    article.Description = description
    ArticlePost(article)
    fmt.Fprintf(w, "Posted some data! [%s]<BR>", content)
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        fmt.Println(err)
        id = 0
    }
    tmpl := template.Must(template.ParseFiles("template/article.html"))
    stuff := ArticleLookup()
    _ = id
    //if len(stuff) >= id && id > 0 {
        //json.NewEncoder(w).Encode(stuff[id-1])
        //tmpl.Execute(w, stuff[id-1])
    //}
    tmpl.Execute(w, stuff[id-1])
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the HomePage!<BR>")
    fmt.Println("")
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
    category := mux.Vars(r)["category"]
    id := mux.Vars(r)["id"]
    fmt.Fprintf(w, "Category: %s<BR>", category)
    fmt.Fprintf(w, "ID: %s<BR>", id)
    //json.NewEncoder(w).Encode(Articles)
}

func CssHandler (w http.ResponseWriter, r *http.Request) {
    css := mux.Vars(r)["css"]
    fmt.Println("Entered CSS")
    path := fmt.Sprintf("./css/%s", css)
    http.ServeFile(w, r, path)
}

func Signin(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    pwd := []byte(r.FormValue("password"))
    hash := hashAndSalt(pwd)
    var credentials Credentials
    credentials.Username = r.FormValue("username")
    credentials.Password = r.FormValue("password")
    success := UserLookup(credentials.Username, credentials.Password)
    fmt.Println(hash)
    if !success {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    u,_ := uuid.NewV4()
    sessionToken := u.String()
    _, err := cache.Do("SETEX", sessionToken, "120", username)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println("Testing ======================================")
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   sessionToken,
        Expires: time.Now().Add(120 * time.Second),
    })
    fmt.Fprintf(w, "Made it!")

}

func Welcome(w http.ResponseWriter, r *http.Request) {
    response, err := checkSession(w, r)
    if err != nil {
        fmt.Fprintf(w, "Failed login check")
        return
    }
    w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
}


func Refresh(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    sessionToken := c.Value

    response, err := cache.Do("GET", sessionToken)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if response == nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    //newSessionToken := uuid.NewV4().String()
    u,_ := uuid.NewV4()
    newSessionToken := u.String()
    _, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s", response))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    _, err = cache.Do("DEL", sessionToken)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   newSessionToken,
        Expires: time.Now().Add(120 * time.Second),
    })
}
