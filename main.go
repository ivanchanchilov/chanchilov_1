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
	db, err := sql.Open("sqlite", "expenses.db?_pragma=foreign_keys(1)")
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
	http.HandleFunc("/register", handlers.RegisterPage)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/logout", handlers.Auth(handlers.LogoutPage))
	http.HandleFunc("/about", handlers.AboutPage)
	http.HandleFunc("/ping", handlers.PingPage)
	http.HandleFunc("/expenses", handlers.Auth(handlers.ExpensesPage))
	http.HandleFunc("/expenses/add", handlers.Auth(handlers.AddExpensePage))
	http.HandleFunc("/expenses/edit", handlers.Auth(handlers.EditExpensePage))
	http.HandleFunc("/expenses/delete", handlers.Auth(handlers.DeleteExpensePage))
	http.HandleFunc("/stats", handlers.Auth(handlers.StatsPage))
	http.HandleFunc("/api/stats/by-category", handlers.Auth(handlers.CategoryStats))
	http.HandleFunc("/api/stats/by-month", handlers.Auth(handlers.MonthStats))
	http.Handle("/static/", http.StripPrefix("/static/", handlers.StaticFiles()))

	http.ListenAndServe(":8080", handlers.Logging(handlers.Recovery(http.DefaultServeMux)))
}
