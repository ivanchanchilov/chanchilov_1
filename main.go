package main

import (
	"database/sql"
	"expenses/internal/handlers"
	"expenses/internal/repository"
	"fmt"
	_ "modernc.org/sqlite"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	repo := repository.ExpenseRepository{DB: db}
	err = repo.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	handlers.Repo = repo

	http.HandleFunc("/", handlers.MainPage)
	http.HandleFunc("/about", handlers.AboutPage)
	http.HandleFunc("/ping", handlers.PingPage)
	http.HandleFunc("/expenses", handlers.ExpensesPage)
	http.HandleFunc("/expenses/add", handlers.AddExpensePage)
	http.HandleFunc("/expenses/edit", handlers.EditExpensePage)
	http.HandleFunc("/expenses/delete", handlers.DeleteExpensePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.ListenAndServe(":8080", nil)
}
