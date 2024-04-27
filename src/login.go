package main

import (
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my login")
}
