package main

import (
	"expenses/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.MainPage)
	http.HandleFunc("/about", handlers.AboutPage)
	http.HandleFunc("/ping", handlers.PingPage)
	http.HandleFunc("/expenses", handlers.ExpensesPage)
	http.HandleFunc("/expenses/add", handlers.AddExpensePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.ListenAndServe(":8080", nil)
}
