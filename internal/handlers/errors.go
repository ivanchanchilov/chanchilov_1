package handlers

import (
	"html/template"
	"net/http"
)

func ErrorPage(w http.ResponseWriter, code int, name string) {
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/"+name+".html")
	if err != nil {
		http.Error(w, http.StatusText(code), code)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	t.ExecuteTemplate(w, "layout", nil)
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	ErrorPage(w, 404, "404")
}

func ServerError(w http.ResponseWriter) {
	ErrorPage(w, 500, "500")
}
