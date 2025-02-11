package main

import (
	"fmt"
	"net/http"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CRSFToken      string
}

var users = map[string]Login{}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)
	fmt.Println("server running on port 8080")
	http.ListenAndServe(":8080", nil)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		er := http.StatusMethodNotAllowed
		http.Error(w, "Invalid method", er)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 4 || len(password) < 8 {
		er := http.StatusNotAcceptable
		http.Error(w, "Invalid username/password", er)
		return
	}

	if _, ok := users[username]; ok {
		er := http.StatusConflict
		http.Error(w, "User already exists", er)
		return
	}

	hashPassword, _ := hashPassword(password)
	users[username] = Login{HashedPassword: hashPassword}

	fmt.Fprintln(w, "User registered successfully")
}

func login(w http.ResponseWriter, r *http.Request) {}

func logout(w http.ResponseWriter, r *http.Request) {}

func protected(w http.ResponseWriter, r *http.Request) {}
