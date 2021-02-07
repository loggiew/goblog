// +build ignore

package main

import (
    "os"
    "fmt"
    "log"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "gopkg.in/ini.v1"
)


type Credentials struct {
    Password string "password"
    Username string "username"
}



func ArticleLookup() []Article {
    cfg, err := ini.Load(".ini")
    if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
    }
    query := `SELECT
    articles.article_id,article_title,article_author,description,content
    FROM articles
    INNER JOIN content
    ON articles.content_id = content.content_id
    INNER JOIN article_meta
    ON articles.article_id = article_meta.article_id;`

    sqlUser := cfg.Section("").Key("user").String()
    sqlPass := cfg.Section("").Key("pass").String()
    connect := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/golang", sqlUser, sqlPass)
    db, err := sql.Open("mysql", connect)

    if err != nil {
        log.Fatal(err)
    }
    rows, err2 := db.Query(query)

    if err2 != nil {
        log.Fatal(err2)
    }
    article := Article{}
    articles := []Article{}
    for rows.Next() {
        var id int
        var title, content, author string
        var description sql.NullString
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

func ArticlePost(article Article) int {
    cfg, err := ini.Load(".ini")
    if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
    }
    sqlUser := cfg.Section("").Key("user").String()
    sqlPass := cfg.Section("").Key("pass").String()
    connect := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/golang", sqlUser, sqlPass)
    db, err := sql.Open("mysql", connect)
    if err != nil {
        log.Fatal(err)
    }
    query1 := `INSERT INTO content (content) VALUES (?);`
    query2 := `INSERT INTO article_meta (article_title, article_author, description) VALUES (?, ?, ?);`

    tx, err := db.Begin()
    handleError(err)

    res, err := tx.Exec(query1, article.Content)
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    last_id1, err := res.LastInsertId()
    handleError(err)

    res, err = tx.Exec(query2, article.Title, article.Auth, article.Description.String)
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    last_id2, err := res.LastInsertId()
    handleError(err)

    query3 := fmt.Sprintf(`
    INSERT INTO articles
    (content_id, article_id)
    VALUES (%d, %d);
    `, last_id1, last_id2)
    res, err = tx.Exec(query3)
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }

    handleError(tx.Commit())
    defer db.Close()
    return 1
}

func UserLookup(username string, password string) (bool){
    cfg, err := ini.Load(".ini")
    if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
    }
    query := fmt.Sprintf(`SELECT
    username,password
    FROM users
    WHERE username = '%s'`, username)

    sqlUser := cfg.Section("").Key("user").String()
    sqlPass := cfg.Section("").Key("pass").String()
    connect := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/golang", sqlUser, sqlPass)
    db, err := sql.Open("mysql", connect)

    if err != nil {
        log.Fatal(err)
    }
    rows, err2 := db.Query(query)

    if err2 != nil {
        log.Fatal(err2)
    }
    var check_credential Credentials
    for rows.Next() {
        var dbusername, dbpassword string
        err = rows.Scan(&dbusername, &dbpassword)
        if err != nil {
            panic(err.Error())
        }
        check_credential.Username = dbusername
        check_credential.Password = dbpassword
    }
    defer db.Close()
    success := comparePasswords(check_credential.Password, []byte(password))
    if success == false {
        return false
    }
    return true
}


func handleError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
