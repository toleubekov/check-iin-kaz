package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/toleubekov/check-iin-kaz/internal/model"
)

type Config struct {
	ServerURL     string
	NumGoroutines int
	NumRequests   int
}

type Person struct {
	Name  string `json:"name"`
	IIN   string `json:"iin"`
	Phone string `json:"phone"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := Config{
		ServerURL:     getEnv("SERVER_URL", "http://server:8080"),
		NumGoroutines: getEnvAsInt("NUM_GOROUTINES", 10),
		NumRequests:   getEnvAsInt("NUM_REQUESTS", 100),
	}

	log.Printf("Starting stress test with %d goroutines, %d requests per goroutine",
		config.NumGoroutines, config.NumRequests)
	log.Printf("Server URL: %s", config.ServerURL)

	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(config.NumGoroutines)

	for i := 0; i < config.NumGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			runStressTest(goroutineID, config)
		}(i)
	}

	wg.Wait()
	log.Println("Stress test completed")
}

func runStressTest(goroutineID int, config Config) {
	for i := 0; i < config.NumRequests; i++ {

		iin := generateRandomIIN(i%10 != 0)

		person := Person{
			Name:  fmt.Sprintf("Test Person %d-%d", goroutineID, i),
			IIN:   iin,
			Phone: fmt.Sprintf("+7%09d", rand.Intn(1000000000)),
		}

		err := createPerson(config.ServerURL, person)
		if err != nil {
			log.Printf("Goroutine %d: Failed to create person: %v", goroutineID, err)
		} else {
			log.Printf("Goroutine %d: Created person with IIN: %s", goroutineID, iin)
		}

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}
}

func createPerson(serverURL string, person Person) error {

	jsonData, err := json.Marshal(person)
	if err != nil {
		return fmt.Errorf("failed to marshal person: %w", err)
	}

	url := fmt.Sprintf("%s/people/info", serverURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response model.PersonResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if !response.Success {
		return fmt.Errorf("failed to create person: %s", response.Errors)
	}

	return nil
}

func generateRandomIIN(valid bool) string {
	year := 1950 + rand.Intn(70)
	month := 1 + rand.Intn(12)
	day := 1 + rand.Intn(28)

	var genderCentury int
	if year >= 2000 {
		genderCentury = 3 + rand.Intn(2)
	} else {
		genderCentury = 1 + rand.Intn(2)
	}

	regNumber := 1000 + rand.Intn(9000)

	iin := fmt.Sprintf("%02d%02d%02d%d%04d", year%100, month, day, genderCentury, regNumber)

	if valid {
		checksum := calculateChecksum(iin)
		iin += strconv.Itoa(checksum)
	} else {
		iin += strconv.Itoa(rand.Intn(10))
	}

	return iin
}

func calculateChecksum(iin11 string) int {
	weights1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

	sum := 0
	for i := 0; i < 11; i++ {
		digit, _ := strconv.Atoi(string(iin11[i]))
		sum += digit * weights1[i]
	}

	controlDigit := sum % 11

	if controlDigit == 10 {
		weights2 := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}

		sum = 0
		for i := 0; i < 11; i++ {
			digit, _ := strconv.Atoi(string(iin11[i]))
			sum += digit * weights2[i]
		}

		controlDigit = sum % 11

		if controlDigit == 10 {
			controlDigit = 0
		}
	}

	return controlDigit
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Could not parse %s=%s as integer, using default %d", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}
