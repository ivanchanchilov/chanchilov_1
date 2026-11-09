package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Headers http.Header
	Body    bytes.Buffer
	Status  int
}

func (response *Response) Header() http.Header {
	return response.Headers
}

func (response *Response) WriteHeader(status int) {
	if response.Status == 0 {
		response.Status = status
	}
}

func (response *Response) Write(data []byte) (int, error) {
	if response.Status == 0 {
		response.Status = 200
	}
	return response.Body.Write(data)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		response := &Response{Headers: make(http.Header)}
		next.ServeHTTP(response, r)
		if response.Status == 0 {
			response.Status = 200
		}
		for name, values := range response.Headers {
			w.Header()[name] = values
		}
		w.WriteHeader(response.Status)
		w.Write(response.Body.Bytes())
		fmt.Println(r.Method, r.URL.Path, response.Status, time.Since(start))
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			problem := recover()
			if problem != nil {
				fmt.Println("Паника:", problem)
				response, ok := w.(*Response)
				if ok {
					response.Headers = make(http.Header)
					response.Body.Reset()
					response.Status = 0
				}
				ServerError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
