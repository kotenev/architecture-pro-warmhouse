package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// TemperatureResponse представляет ответ от API
type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	// Обрабатываем запросы как по query-параметру, так и по ID в пути
	http.HandleFunc("/temperature", temperatureHandler)
	http.HandleFunc("/temperature/", temperatureHandler)

	log.Println("Temperature API starting on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем sensorID из пути (например, /temperature/1)
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	sensorID := ""
	if len(pathParts) > 1 {
		sensorID = pathParts[1]
	}

	// Извлекаем location из query-параметра (например, ?location=Living+Room)
	location := r.URL.Query().Get("location")

	// Реализуем логику из задания
	// Если location не указан, определяем его по sensorID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// Если sensorID не указан, определяем его по location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	// Генерируем случайные данные
	randomTemp := 17.0 + rand.Float64()*(10.0) // Температура от 18.0 до 28.0
	status := "active"
	if rand.Intn(10) > 8 { // 20% шанс, что датчик неактивен
		status = "inactive"
	}

	// Формируем ответ
	response := TemperatureResponse{
		Value:       float64(int(randomTemp*10)) / 10, // Округляем до 1 знака после запятой
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      status,
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: fmt.Sprintf("Real-time temperature data for %s", location),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
