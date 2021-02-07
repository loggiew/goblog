// +build ignore

package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)


func main() {
    initCache()
    router := mux.NewRouter().StrictSlash(false)
    router.HandleFunc("/css/{css}", CssHandler)
    router.HandleFunc("/", homePage).Methods("GET")
    router.HandleFunc("/test", test)
    router.HandleFunc("/login", Signin).Methods("POST")
    router.HandleFunc("/welcome", Welcome)
    router.HandleFunc("/refresh", Refresh)
    router.HandleFunc("/createpost", CreatePost).Methods("GET")
    router.HandleFunc("/articles", AllArticlesGet).Methods("GET")
    router.HandleFunc("/articles", AllArticlesPost).Methods("POST")
    router.HandleFunc("/articles/{id:[0-9]+}", ArticleHandler).Methods("GET")
    router.HandleFunc("/articles/{category}/", ArticlesCategoryHandler).Methods("GET")
    router.HandleFunc("/product/{category}/{id:[0-9]+}", ProductHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", removeTrailingSlash(router)))
}

