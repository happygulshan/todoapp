package main

import (
	"log"
	"net/http"
	"todoapp/db"
	"todoapp/handlers"

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

	r.HandleFunc("/signup", h.CreateUser).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/logout", h.Logout).Methods("POST")

	// task related job here
	r.HandleFunc("/createtask", h.CreateTask).Methods("POST")
	r.HandleFunc("/getalltask", h.GetAllTasks).Methods("GET")
	r.HandleFunc("/updatetask/{id}", h.UpdateTask).Methods("PATCH")
	r.HandleFunc("/deletetask/{id}", h.DeleteTask).Methods("DELETE")

	log.Println("Server starting on :8080")

	http.ListenAndServe(":8080", r)
}
