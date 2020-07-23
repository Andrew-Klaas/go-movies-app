package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Andrew-Klaas/go-movies-app/movies"
	"github.com/Andrew-Klaas/go-movies-app/users"
)

func init() {
	users.Users["test@test.com"] = users.User{"test@test.com", []byte("123"), "a", "k", "user"}
}

type appInfo struct {
	Version int    `json: "Version"`
	Owner   string `json: "Owner"`
}

func main() {
	//Login
	http.HandleFunc("/", users.Index)
	http.HandleFunc("/signup", users.Signup)
	http.HandleFunc("/login", users.Login)
	http.HandleFunc("/logout", users.Logout)
	http.HandleFunc("/version", Version)
	//moviestore
	http.HandleFunc("/movies", movies.MovieStore)
	http.HandleFunc("/moviestore", movies.MovieStore)
	http.HandleFunc("/movies/show", movies.Show)
	http.HandleFunc("/movies/create", movies.Create)
	http.HandleFunc("/movies/create/process", movies.CreateProcess)
	http.HandleFunc("/movies/update", movies.Update)
	http.HandleFunc("/movies/update/process", movies.UpdateProcess)
	http.HandleFunc("/movies/delete/process", movies.DeleteProcess)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	//Server
	http.ListenAndServe(":8080", nil)
}

//Version ...
func Version(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ai := appInfo{
		Version: 2,
		Owner:   "andrew",
	}
	err := json.NewEncoder(w).Encode(ai)
	if err != nil {
		log.Println(err)
	}
}
