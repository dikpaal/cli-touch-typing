package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	// you have to load the .env first
	_ = godotenv.Load()

	endpoint := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=Iasi&days=1&aqi=no&alerts=no", os.Getenv("API_KEY"))

	res, err := http.Get(endpoint)

	if err != nil {
		panic(err)
	}

	// close the response body when done

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Error: Unable to fetch weather data")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var weather Weather

	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	fmt.Println(weather)
}
