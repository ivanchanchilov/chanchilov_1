package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"html/template"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDKey contextKey = "userID"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		userID, err := Repo.GetSessionUser(cookie.Value, time.Now().Unix())
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if err != nil {
			ServerError(w)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("web/templates/layout.html", "web/templates/register.html")
		if err != nil {
			ServerError(w)
			return
		}
		t.ExecuteTemplate(w, "layout", nil)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", 400)
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")
	if email == "" || password == "" || len(password) > 72 {
		AuthForm(w, r, "register", "Введите email и пароль длиной до 72 байт.", 400)
		return
	}
	_, err = Repo.GetUser(email)
	if err == nil {
		AuthForm(w, r, "register", "Этот email уже зарегистрирован. Войдите в аккаунт.", 400)
		return
	}
	if err != sql.ErrNoRows {
		ServerError(w)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ServerError(w)
		return
	}
	err = Repo.CreateUser(email, string(hash))
	if err != nil {
		AuthForm(w, r, "register", "Не удалось создать аккаунт. Попробуйте другой email.", 400)
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("web/templates/layout.html", "web/templates/login.html")
		if err != nil {
			ServerError(w)
			return
		}
		t.ExecuteTemplate(w, "layout", nil)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", 400)
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	user, err := Repo.GetUser(email)
	if err == sql.ErrNoRows {
		AuthForm(w, r, "login", "Не удалось войти. Проверьте email и пароль.", 401)
		return
	}
	if err != nil {
		ServerError(w)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(r.FormValue("password")))
	if err != nil {
		AuthForm(w, r, "login", "Не удалось войти. Проверьте email и пароль.", 401)
		return
	}
	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		ServerError(w)
		return
	}
	sessionID := hex.EncodeToString(bytes)
	expires := time.Now().Add(24 * time.Hour)
	err = Repo.CreateSession(sessionID, user.ID, expires.Unix())
	if err != nil {
		ServerError(w)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
	http.Redirect(w, r, "/expenses", 302)
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	cookie, err := r.Cookie("session")
	if err == nil {
		err = Repo.DeleteSession(cookie.Value)
		if err != nil {
			ServerError(w)
			return
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
	http.Redirect(w, r, "/login", 302)
}

type AuthData struct {
	Email string
	Error string
}

func AuthForm(w http.ResponseWriter, r *http.Request, name string, message string, status int) {
	data := AuthData{Email: r.FormValue("email"), Error: message}
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/"+name+".html")
	if err != nil {
		ServerError(w)
		return
	}
	w.WriteHeader(status)
	t.ExecuteTemplate(w, "layout", data)
}
