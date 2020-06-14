package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// SunsetAPI ...
type SunsetAPI struct {
	Results TimeInformations `json:"results"`
	Status  string           `json:"status"`
}

// TimeInformations ...
type TimeInformations struct {
	Sunrise   string `json:"sunrise"`
	Sunset    string `json:"sunset"`
	DayLength string `json:"day_length"`
}

type WeatherAPI struct {
	Coord   Coord     `json:"coord"`
	Weather []Weather `json:"weather"`
	Main    All       `json:"main"`
	Wind    Wind      `json:"wind"`
	Clouds  Clouds    `json:"clouds"`
	Dt      int       `json:"dt"`
	Sys     Sys       `json:"sys"`
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Cod     int       `json:"cod"`
}
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
type All struct {
	Temp     float64 `json:"temp"`
	Pressure float64 `json:"pressure"`
	Humidity int     `json:"humidity"`
	TempMin  float64 `json:"temp_min"`
	TempMax  float64 `json:"temp_max"`
}
type Wind struct {
	Speed float64 `json:"speed"`
}
type Clouds struct {
	All int `json:"all"`
}
type Sys struct {
	Message float64 `json:"message"`
	Country string  `json:"country"`
}

// APIResponse CCR API Main response
type APIResponse struct {
	TimeInformations    *TimeInformations `json:"time"`
	WeatherInformations *WeatherAPI       `json:"wether"`
}

func sunset(lat, long string) (*TimeInformations, error) {
	url := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%s&lng=%s", lat, long)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sunsetAPI := &SunsetAPI{}

	if err = json.Unmarshal(body, sunsetAPI); err != nil {
		return nil, err
	}

	ss := sunsetAPI.Results

	return &ss, nil

}

func weather(lat, long string) (*WeatherAPI, error) {
	appid := os.Getenv("OPEN_WETHER_APPID")
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, long, appid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	weatherAPI := &WeatherAPI{}

	if err = json.Unmarshal(body, weatherAPI); err != nil {
		return nil, err
	}

	wa := weatherAPI

	return wa, nil
}

func glue(lat, long string) (*APIResponse, error) {
	timeInfo, err := sunset(lat, long)
	if err != nil {
		return nil, err
	}

	weatherInfo, err := weather(lat, long)
	if err != nil {
		return nil, err
	}

	return &APIResponse{
		TimeInformations:    timeInfo,
		WeatherInformations: weatherInfo,
	}, nil

}

func all(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	lat := query.Get("lat")
	long := query.Get("long")

	apiResp, err := glue(lat, long)
	if err != nil {

	}

	body, err := json.Marshal(apiResp)
	if err != nil {

	}
	w.Write(body)
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/", all).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
