package handlers

import (
	"encoding/json"
	"expenses/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func StatsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	year := time.Now().Year()
	var err error
	if r.URL.Query().Get("year") != "" {
		year, err = strconv.Atoi(r.URL.Query().Get("year"))
		if err != nil || year < 1 || year > 9999 {
			http.Error(w, "Неверный год", 400)
			return
		}
	}
	data := models.StatsPageData{
		From: r.URL.Query().Get("from"),
		To:   r.URL.Query().Get("to"),
		Year: year,
	}
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/stats.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", 500)
		return
	}
	t.ExecuteTemplate(w, "layout", data)
}

func CategoryStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	userID := r.Context().Value(userIDKey).(int)
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	stats, err := Repo.GetCategoryStats(userID, from, to)
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func MonthStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	year := time.Now().Year()
	var err error
	if r.URL.Query().Get("year") != "" {
		year, err = strconv.Atoi(r.URL.Query().Get("year"))
		if err != nil || year < 1 || year > 9999 {
			http.Error(w, "Неверный год", 400)
			return
		}
	}
	userID := r.Context().Value(userIDKey).(int)
	stats, err := Repo.GetMonthStats(userID, year)
	if err != nil {
		http.Error(w, "Ошибка базы данных", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
