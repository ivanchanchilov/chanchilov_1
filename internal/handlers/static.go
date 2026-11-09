package handlers

import (
	"net/http"
	"os"
)

func StaticFiles() http.Handler {
	files := http.FileServer(http.Dir("web/static"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := http.Dir("web/static").Open(r.URL.Path)
		if os.IsNotExist(err) {
			NotFoundPage(w, r)
			return
		}
		if err != nil {
			ServerError(w)
			return
		}
		file.Close()
		files.ServeHTTP(w, r)
	})
}
