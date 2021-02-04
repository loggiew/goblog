// +build ignore

package main

import (
    "fmt"
    "log"
    "strings"
    "strconv"
    "net/http"
    "database/sql"
    "html/template"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-redis/redis"
//    "encoding/json"
    "github.com/gorilla/mux"
    //"github.com/go-session/redis"
    //"github.com/go-session/session"
)


type Article struct {
    Id int `json:"Id"`
    Title string `json:"Title"`
    Auth string `json:"Author"`
    Description string `json:"Description"`
    Content string `json:"Content"`
}

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!<BR>")
    fmt.Fprintf(w, "Hi there, I love %s!<BR>", r.URL.Path[1:])
    fmt.Println("Endpoint Hit: homePage")
}

func LoginHandler(w http.ResponseWriter, r *http.Request){

}

func AllArticlesGet(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    tmpl := template.Must(template.ParseFiles("template/articles.html"))
    stuff := mysql_connect()
    tmpl.Execute(w, stuff)
    //json.NewEncoder(w).Encode(Articles)
}

func AllArticlesPost(w http.ResponseWriter, r *http.Request){
    //content, err := strconv.Atoi(mux.Vars(r)["content"])
    //content := mux.Vars(r)["content"]
    content := r.FormValue("content")
    //if err != nil {
    //    fmt.Println(err)
    //}
    fmt.Println("Entered Articles Post")
    fmt.Fprintf(w, "Posted some data! [%s]<BR>", content)
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        fmt.Println(err)
        id = 0
    }
    tmpl := template.Must(template.ParseFiles("template/article.html"))
    stuff := mysql_connect()
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

func removeTrailingSlash(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/html")
	if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
        }
        next.ServeHTTP(w, r)
    })
}


func mysql_connect() []Article {
    db, err := sql.Open("mysql", "someuser:somepass@tcp(127.0.0.1:3306)/golang")

    if err != nil {
        log.Fatal(err)
    }
    rows,err2 := db.Query(`SELECT 
        articles.article_id,article_title,article_author,description,content 
        FROM articles 
        INNER JOIN content 
        ON articles.content_id = content.content_id 
        INNER JOIN article_meta 
        ON articles.article_id = article_meta.article_id;`)

    if err2 != nil {
        log.Fatal(err2)
    }
    article := Article{}
    articles := []Article{}
    for rows.Next() {
        var id int
        var title, description, content, author string
        err = rows.Scan(&id, &title, &author, &description, &content)
        if err != nil {
            panic(err.Error())
        }
        article.Id = id
        article.Title = title
        article.Auth = author
        article.Description = description
        article.Content = content
        articles = append(articles, article)
    }
    defer db.Close()
    return articles
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	pong, err := client.Ping(client.Context()).Result()
        fmt.Println(pong, err)
    router := mux.NewRouter().StrictSlash(false)
    router.HandleFunc("/css/{css}", CssHandler)
    router.HandleFunc("/", homePage).Methods("GET")
    router.HandleFunc("/login", LoginHandler).Methods("GET")
    router.HandleFunc("/articles", AllArticlesGet).Methods("GET")
    router.HandleFunc("/articles", AllArticlesPost).Methods("POST")
    router.HandleFunc("/articles/{id:[0-9]+}", ArticleHandler).Methods("GET")
    router.HandleFunc("/articles/{category}/", ArticlesCategoryHandler).Methods("GET")
    router.HandleFunc("/product/{category}/{id:[0-9]+}", ProductHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", removeTrailingSlash(router)))
    //log.Fatal(http.ListenAndServe(":8080", router))
}
/*

        */
