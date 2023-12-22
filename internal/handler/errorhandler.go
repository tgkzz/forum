package handler

import (
	"fmt"
	"forum/internal/model"
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	tmpl, err := template.ParseFiles("template/html/error.html")
	if err != nil {
		text := fmt.Sprintf("Error 500\n Oppss! %s", http.StatusText(http.StatusInternalServerError))
		http.Error(w, text, http.StatusInternalServerError)
		return
	}
	res := &model.Err{Text_err: http.StatusText(code), Code_err: code}
	err = tmpl.Execute(w, &res)
	if err != nil {
		text := fmt.Sprintf("Error 500\n Oppss! %s", http.StatusText(http.StatusInternalServerError))
		http.Error(w, text, http.StatusInternalServerError)
		return
	}
}

func ClientErrorHandler(tmpl *template.Template, w http.ResponseWriter, cErr error, statuscode int) {
	type ClientError struct {
		ErrorText string
	}

	w.WriteHeader(statuscode)

	err := tmpl.Execute(w, ClientError{
		ErrorText: cErr.Error(),
	})
	if err != nil {
		log.Print(err)
		ErrorHandler(w, http.StatusInternalServerError)
	}
}
