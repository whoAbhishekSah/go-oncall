package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/oncall?sslmode=disable"
	}
	
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	initDB()
	
	r := mux.NewRouter()
	
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users", getUsersHandler).Methods("GET")
	r.HandleFunc("/teams", createTeamHandler).Methods("POST")
	r.HandleFunc("/teams", getTeamsHandler).Methods("GET")
	r.HandleFunc("/schedules", createScheduleHandler).Methods("POST")
	r.HandleFunc("/schedules", getSchedulesHandler).Methods("GET")
	
	go scheduleChecker()
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
