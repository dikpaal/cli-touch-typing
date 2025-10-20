package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

	fmt.Println(string(body))

}
