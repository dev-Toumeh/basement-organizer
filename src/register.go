package main

import (
    "html/template"
    "net/http"
    "log"
)

type registerPage struct {

  Title string
}

func register(response http.ResponseWriter, request *http.Request) {
    if request.Method != "GET" {
        http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    tmpl, err := template.ParseFiles("./register.html")
    if err != nil {
        http.Error(response, "Error loading template", http.StatusInternalServerError)
        return
    }
    data := registerPage{
        Title: "Register",
    }
    err = tmpl.Execute(response, data)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(response, "Error rendering page", http.StatusInternalServerError)
    }
}
