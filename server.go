package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

const openWeatherAPIKey = ""

type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type WeatherResponse struct {
	Date        string  `json:"date"`
	City        string  `json:"city"`
	Unit        string  `json:"unit"`
	Temperature float64 `json:"temperature"`
}

func main() {
	http.HandleFunc("/weather", weatherCurrent)
	fmt.Println("Server started. Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func weatherCurrent(w http.ResponseWriter, r *http.Request) {

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is missing", http.StatusBadRequest)
		return
	}

	weatherTemp, err := getWeather(strings.ToLower(city))
	if err != nil {
		http.Error(w, "Failed to get weather data", http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{
		Date:        time.Now().Format("2006-01-02T15:04:05"),
		City:        city,
		Unit:        "celsius",
		Temperature: weatherTemp,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getWeather(city string) (float64, error) {

	// openWeatherAPIKey := os.Getenv("key")

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, openWeatherAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to fetch weather data")
	}

	var weatherData WeatherData
	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		return 0, err
	}
	temperature := math.Round(weatherData.Main.Temp - 273.15)
	return temperature, nil
}
