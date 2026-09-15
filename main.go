package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/about", aboutPage)
	http.HandleFunc("/ping", pingPage)

	http.ListenAndServe(":8080", nil)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "123123")
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "321321321")
}

func pingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintln(w, "pong")
		return
	}
	http.Error(w, "Method now allowed", 405)
}
