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
		http.NotFound(w, r)
		return
	}

	ExpensesPage(w, r)
}

func AboutPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/about.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", 500)
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
	expenses, err := Repo.GetAll(categoryID, from, to)
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
		return
	}

	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/expenses.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", 500)
		return
	}

	categories, err := Repo.GetCategories()
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
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
	if r.Method == "GET" {
		t, err := template.ParseFiles("web/templates/layout.html", "web/templates/add.html")
		if err != nil {
			http.Error(w, "Ошибка шаблона", 500)
			return
		}

		categories, err := Repo.GetCategories()
		if err != nil {
			http.Error(w, "Ошибка базы данных", 500)
			return
		}
		data := models.PageData{Categories: categories}
		t.ExecuteTemplate(w, "layout", data)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка формы", 400)
			return
		}

		amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil {
			http.Error(w, "Введите сумму числом", 400)
			return
		}

		categoryID, err := strconv.Atoi(r.FormValue("category"))
		if err != nil {
			http.Error(w, "Неверная категория", 400)
			return
		}

		expense := models.Expense{
			CategoryID:  categoryID,
			Amount:      amount,
			Description: r.FormValue("description"),
			Date:        r.FormValue("date"),
		}

		err = Repo.Create(expense)
		if err != nil {
			http.Error(w, "Ошибка базы данных", 500)
			return
		}
		http.Redirect(w, r, "/expenses", 302)
		return
	}

	http.Error(w, "Method not allowed", 405)
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

	expense, err := Repo.GetByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
		return
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("web/templates/layout.html", "web/templates/edit.html")
		if err != nil {
			http.Error(w, "Ошибка шаблона", 500)
			return
		}
		categories, err := Repo.GetCategories()
		if err != nil {
			http.Error(w, "Ошибка базы данных", 500)
			return
		}
		data := models.PageData{Expense: expense, Categories: categories}
		t.ExecuteTemplate(w, "layout", data)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", 400)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Введите сумму числом", 400)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category"))
	if err != nil {
		http.Error(w, "Неверная категория", 400)
		return
	}
	expense.CategoryID = categoryID
	expense.Amount = amount
	expense.Description = r.FormValue("description")
	expense.Date = r.FormValue("date")
	err = Repo.Update(expense)
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
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

	expense, err := Repo.GetByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
		return
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("web/templates/layout.html", "web/templates/delete.html")
		if err != nil {
			http.Error(w, "Ошибка шаблона", 500)
			return
		}
		t.ExecuteTemplate(w, "layout", expense)
		return
	}

	err = Repo.Delete(id)
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
		return
	}

	http.Redirect(w, r, "/expenses", 302)
}
