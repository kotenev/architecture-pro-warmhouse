package main

import (
	"database/sql"
	"encoding/json"
	_ "fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type TelemetryData struct {
	SensorID  int       `json:"sensor_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

var db *sql.DB

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Подключение к БД
	connStr := os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", connStr)
	failOnError(err, "Failed to connect to database")
	defer db.Close()

	// Подключение к RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("telemetry", true, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	// Горутина для прослушивания RabbitMQ
	go func() {
		for d := range msgs {
			var data TelemetryData
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Error decoding telemetry message: %s", err)
				continue
			}
			log.Printf("Received telemetry: %+v", data)
			// Сохранение в БД
			_, err = db.Exec(
				"INSERT INTO telemetry_data (sensor_id, value, timestamp) VALUES ($1, $2, $3)",
				data.SensorID, data.Value, data.Timestamp)
			if err != nil {
				log.Printf("Failed to save telemetry to DB: %v", err)
			}
		}
	}()

	// Этот эндпоинт будет имитировать старый temperature-api для монолита
	http.HandleFunc("/temperature", getTemperature)
	log.Println("Telemetry Service starting on port 8083...")
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func getTemperature(w http.ResponseWriter, r *http.Request) {
	// Для простоты MVP, мы будем возвращать последние данные для сенсора 1
	var value float64
	var timestamp time.Time
	err := db.QueryRow("SELECT value, timestamp FROM telemetry_data WHERE sensor_id = 1 ORDER BY timestamp DESC LIMIT 1").Scan(&value, &timestamp)
	if err != nil {
		// Если данных нет, вернем случайное значение, чтобы монолит не падал
		value = 20.0 + rand.Float64()*5.0
		timestamp = time.Now()
	}

	// Структура ответа, которую ожидает монолит
	response := struct {
		Value       float64   `json:"value"`
		Unit        string    `json:"unit"`
		Timestamp   time.Time `json:"timestamp"`
		Location    string    `json:"location"`
		Status      string    `json:"status"`
		SensorID    string    `json:"sensor_id"`
		SensorType  string    `json:"sensor_type"`
		Description string    `json:"description"`
	}{
		Value:       value,
		Unit:        "°C",
		Timestamp:   timestamp,
		Location:    "Living Room",
		Status:      "active",
		SensorID:    "1",
		SensorType:  "temperature",
		Description: "Real-time from Telemetry Service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
