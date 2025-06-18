package main

import (
	"log"
	"net/http"
	"todoapp/db"
	"todoapp/handlers"

	"github.com/gorilla/mux"
)

func main() {

	database := db.InitDB()
	defer database.Close()

	db.RunMigrations(database)

	h := handlers.Handler{DB: database}
	r := mux.NewRouter()

	// authMiddleware := middleware.AuthMiddleware(database)

	r.HandleFunc("/signup", h.CreateUser).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/logout", h.Logout).Methods("POST")
	r.HandleFunc("/task", h.CreateTask).Methods("POST")


	// Protect routes, will implement below later
	// r.Handle("/task", authMiddleware(http.HandlerFunc(h.CreateTask))).Methods("POST")

	// r.HandleFunc("/task", h.Getalltasks).Method("GET")
	// r.HandleFunc("/task/{id}", h.Getalltasks).Method("GET")

	log.Println("Server starting on :8080")

	http.ListenAndServe(":8080", r)
}
