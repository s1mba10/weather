package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TempC float64 `json:"temp_c"`
				TimeEpoch int64 `json:"time_epoch"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain int64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Moscow"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=ad3cf5ffcee8449188b103358240209&q=" + q + "&aqi=no")
	if err != nil {
		panic(err)
	} 
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API is not available")
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

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	
	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, 
	location.Country, 
	current.TempC, 
	current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0fC, %d%% chance of rain, %s\n", 
		date.Format("15:04"), 		// Время по часам
		hour.TempC, 				// Температура
		hour.ChanceOfRain, 			// Вероятность дождя
		hour.Condition.Text,)		// Состояние погоды

		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
		
}
