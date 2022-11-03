package main

import (
	"log"
	"net/http"

	"example.com/basicWebApp/handlers"
)

func main() {
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/registerauth", handlers.RegisterAuthHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/loginauth", handlers.LoginAuthHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

