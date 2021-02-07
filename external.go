// +build ignore

package main

import (
    "fmt"
    "errors"
    "log"
    "strings"
    "net/http"
    "github.com/gomodule/redigo/redis"
    "golang.org/x/crypto/bcrypt"
)

var cache redis.Conn

func initCache(){
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	cache = conn
}

func checkSession(w http.ResponseWriter, r *http.Request) (interface {}, error){
    var response interface {}
    c, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            //w.WriteHeader(http.StatusUnauthorized)
            return response, errors.New("Unauthorized")
        }
        //w.WriteHeader(http.StatusBadRequest)
        return response,  errors.New("Unauthorized")
    }

    sessionToken := c.Value

    response, err = cache.Do("GET", sessionToken)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return "", errors.New("Internal Error")
    }
    if response == nil {
        //w.WriteHeader(http.StatusUnauthorized)
        return "", errors.New("Unauthorized")
    }
    fmt.Println(response)
    return response, nil
}

func hashAndSalt(pwd []byte) string {
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }
    return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
    byteHash := []byte(hashedPwd)
    err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
    if err != nil {
        log.Println(err)
        return false
    }
    return true
}

func removeTrailingSlash(next http.Handler) http.Handler {
    // Remove trailing slash
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/html")
        response, err := checkSession(w, r)
        if err == nil {
            w.Write([]byte(fmt.Sprintf("Hello, %s\n<BR>", response)))
        }
	if r.URL.Path != "/" {
            r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
        }
        next.ServeHTTP(w, r)
    })
}
