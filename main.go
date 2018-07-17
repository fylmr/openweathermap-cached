package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"github.com/fylmr/openweathermap-cached/fileCache"
)

type Config struct {
	WEATHER_HOST string
	API_KEY      string
	SERVER_IP    string
	SERVER_PORT  string
	CACHE_TIME   int    // In minutes
	LANGUAGE     string // Two-letter code like "ru"
	UNITS        string // like "metric"
}

var config Config

func main() {
	fmt.Println("Weather is loading...")

	readConfigFile()

	r := mux.NewRouter()
	r.HandleFunc("/{lat}/{lon}", LatLon)

	srv := &http.Server{
		Addr:         config.SERVER_IP + ":" + config.SERVER_PORT,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      r,
	}

	log.Println("Listening at", config.SERVER_IP + ":" + config.SERVER_PORT)

	log.Fatal(srv.ListenAndServe()) // Server launching
}

func LatLon(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	lat := vars["lat"]
	lon := vars["lon"]

	weather, err := fileCache.GetCachedWeather(lat, lon, config.CACHE_TIME)

	if err == nil {
		showJson(writer, weather)
		return
	}

	weather, err = getWeatherFromInternet(lat, lon)
	if err == nil {
		fileCache.SaveToCache(lat, lon, weather)
	}

	showJson(writer, weather)
}

func showJson(writer http.ResponseWriter, weather string) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(weather))
}

func getWeatherFromInternet(lat string, lon string) (weather string, err error) {
	log.Println("Getting weather from the server...")

	url := config.WEATHER_HOST +
		`?lat=` + lat + `&lon=` + lon +
		`&APPID=` + config.API_KEY +
		`&lang=` + config.LANGUAGE + `&units=` + config.UNITS

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("Got weather from the server.")

	weather = string(bs)
	return weather, nil
}

func readConfigFile() {
	log.Println("Reading config file...")

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}

	json.Unmarshal(file, &config)

	log.Println("Config ready.")
}
