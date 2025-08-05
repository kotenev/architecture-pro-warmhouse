package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Device struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	connStr := os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/devices", devicesHandler)
	log.Println("Device Service starting on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func devicesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getDevices(w, r)
	case http.MethodPost:
		createDevice(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, location, type, created_at FROM devices")
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	devices := []Device{}
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Location, &d.Type, &d.CreatedAt); err != nil {
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			return
		}
		devices = append(devices, d)
	}
	json.NewEncoder(w).Encode(devices)
}

func createDevice(w http.ResponseWriter, r *http.Request) {
	var d Device
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO devices (name, location, type) VALUES ($1, $2, $3) RETURNING id, created_at",
		d.Name, d.Location, d.Type,
	).Scan(&d.ID, &d.CreatedAt)

	if err != nil {
		http.Error(w, fmt.Sprintf("Database insert error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}
