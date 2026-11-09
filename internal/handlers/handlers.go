package handlers

import (
	"database/sql"
	"expenses/internal/models"
	"expenses/internal/repository"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var Repo repository.ExpenseRepository

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFoundPage(w, r)
		return
	}

	Auth(ExpensesPage)(w, r)
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/about.html")
	if err != nil {
		ServerError(w)
		return
	}

	t.ExecuteTemplate(w, "layout", nil)
}

func PingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintln(w, "pong")
		return
	}

	http.Error(w, "Method not allowed", 405)
}

func ExpensesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	categoryID := 0
	var err error
	category := r.URL.Query().Get("category")
	if category != "" {
		categoryID, err = strconv.Atoi(category)
		if err != nil {
			http.Error(w, "Неверная категория", 400)
			return
		}
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	expenses, err := Repo.GetAll(r.Context().Value(userIDKey).(int), categoryID, from, to)
	if err != nil {
		ServerError(w)
		return
	}

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/expenses.html")
	if err != nil {
		ServerError(w)
		return
	}

	categories, err := Repo.GetCategories()
	if err != nil {
		ServerError(w)
		return
	}
	data := models.PageData{
		Expenses:   expenses,
		Categories: categories,
		CategoryID: categoryID,
		From:       from,
		To:         to,
	}
	t.ExecuteTemplate(w, "layout", data)
}

func AddExpensePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	categories, err := Repo.GetCategories()
	if err != nil {
		ServerError(w)
		return
	}
	if r.Method == "GET" {
		ExpenseForm(w, r, "add", models.Expense{}, categories, "")
		return
	}
	err = r.ParseForm()
	if err != nil {
		ExpenseForm(w, r, "add", models.Expense{}, categories, "Не удалось прочитать форму.")
		return
	}
	expense, message := ReadExpense(r, categories)
	if message != "" {
		ExpenseForm(w, r, "add", expense, categories, message)
		return
	}
	expense.UserID = r.Context().Value(userIDKey).(int)
	err = Repo.Create(expense)
	if err != nil {
		ServerError(w)
		return
	}
	http.Redirect(w, r, "/expenses", 302)
}

func EditExpensePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Неверный id", 400)
		return
	}
	userID := r.Context().Value(userIDKey).(int)
	expense, err := Repo.GetByID(id, userID)
	if err == sql.ErrNoRows {
		NotFoundPage(w, r)
		return
	}
	if err != nil {
		ServerError(w)
		return
	}
	categories, err := Repo.GetCategories()
	if err != nil {
		ServerError(w)
		return
	}
	if r.Method == "GET" {
		ExpenseForm(w, r, "edit", expense, categories, "")
		return
	}
	err = r.ParseForm()
	if err != nil {
		ExpenseForm(w, r, "edit", expense, categories, "Не удалось прочитать форму.")
		return
	}
	updated, message := ReadExpense(r, categories)
	updated.ID = id
	updated.UserID = userID
	if message != "" {
		ExpenseForm(w, r, "edit", updated, categories, message)
		return
	}
	err = Repo.Update(updated)
	if err != nil {
		ServerError(w)
		return
	}
	http.Redirect(w, r, "/expenses", 302)
}

func DeleteExpensePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Неверный id", 400)
		return
	}

	expense, err := Repo.GetByID(id, r.Context().Value(userIDKey).(int))
	if err == sql.ErrNoRows {
		NotFoundPage(w, r)
		return
	}
	if err != nil {
		ServerError(w)
		return
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("web/templates/layout.html", "web/templates/delete.html")
		if err != nil {
			ServerError(w)
			return
		}
		t.ExecuteTemplate(w, "layout", expense)
		return
	}

	err = Repo.Delete(id, r.Context().Value(userIDKey).(int))
	if err != nil {
		ServerError(w)
		return
	}

	http.Redirect(w, r, "/expenses", 302)
}
