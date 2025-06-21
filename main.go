package main

import (
	"log"
	"net/http"
	"todoapp/db"
	"todoapp/handlers"
	"todoapp/middleware"

	"github.com/gorilla/mux"
)

func main() {

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	defer database.Close()

	db.RunMigrations(database)

	h := handlers.Handler{DB: database}
	r := mux.NewRouter()
	authMiddleware := middleware.AuthMiddleware(database)

	r.HandleFunc("/signup", h.CreateUser).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/logout", h.Logout).Methods("POST")

	// task related job here
	r.Handle("/createtask", authMiddleware(http.HandlerFunc(h.CreateTask))).Methods("POST")
	r.Handle("/getalltask", authMiddleware(http.HandlerFunc(h.GetAllTasks))).Methods("GET")
	r.Handle("/updatetask/{id}", authMiddleware(http.HandlerFunc(h.UpdateTask))).Methods("PATCH")
	r.Handle("/deletetask/{id}", authMiddleware(http.HandlerFunc(h.DeleteTask))).Methods("DELETE")

	log.Println("Server starting on :8080")

	http.ListenAndServe(":8080", r)
}
