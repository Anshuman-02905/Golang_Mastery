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

type APIConfig struct {
	APIKey string
}

type WeatherResponse struct {
	CityName string `json:"name"`
	Main     struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

type Message struct {
	Msg string `json:"message"`
}

func loadConfig() (APIConfig, error) {
	var config APIConfig

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return config, err
	}

	key, exists := os.LookupEnv("OPENWEATHER_API_KEY")
	if !exists {
		log.Fatal("Missing API key in environment variables")
		return config, errors.New("missing API key")
	}

	config = APIConfig{APIKey: key}
	return config, nil
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/weather/", weatherHandler)

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Message{Msg: "Hi, welcome to Weather API BOT!"}
	json.NewEncoder(w).Encode(response)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "City name is required", http.StatusBadRequest)
		return
	}

	city := pathParts[2]
	data, err := getTemperature(city)
	if err != nil {
		http.Error(w, "Failed to fetch city data", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func getTemperature(city string) (*WeatherResponse, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load API key: %v", err)
	}

	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s", config.APIKey, city)
	res, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var weatherData WeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &weatherData, nil
}
