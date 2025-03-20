package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// APIConfig holds the API key configuration
type APIConfig struct {
	APIKey string
}

// WeatherResponse struct maps the JSON response from OpenWeather API
// Contains city name and main temperature details
type WeatherResponse struct {
	CityName string `json:"name"`
	Main     struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

// Message struct is used to send simple JSON responses
// Example: Welcome message response
type Message struct {
	Msg string `json:"message"`
}

// loadConfig reads the API key from the .env file
// Returns APIConfig struct containing the key
func loadConfig() (APIConfig, error) {
	var config APIConfig

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return config, err
	}

	// Retrieve API key from environment variables
	key, exists := os.LookupEnv("OPENWEATHER_API_KEY")
	if !exists {
		log.Fatal("Missing API key in environment variables")
		return config, errors.New("missing API key")
	}

	config = APIConfig{APIKey: key}
	return config, nil
}

func main() {
	// Register HTTP handlers
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/weather/", weatherHandler)

	// Start the server on port 8080
	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// helloHandler responds with a welcome message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Message{Msg: "Hi, welcome to Weather API BOT!"}
	json.NewEncoder(w).Encode(response)
}

// weatherHandler processes weather requests for a city
// Extracts city name from the URL and fetches weather data
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	// Extract city name from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "City name is required", http.StatusBadRequest)
		return
	}

	city := pathParts[2]

	// Fetch temperature data
	data, err := getTemperature(city)
	if err != nil {
		http.Error(w, "Failed to fetch city data", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	json.NewEncoder(w).Encode(data)
}

// getTemperature fetches weather data from OpenWeather API
// Returns a WeatherResponse struct containing temperature details
func getTemperature(city string) (*WeatherResponse, error) {
	// Load API key from configuration
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load API key: %v", err)
	}

	// Construct API URL
	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s", config.APIKey, city)

	// Send GET request to OpenWeather API
	res, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse JSON response
	var weatherData WeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &weatherData, nil
}
