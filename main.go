package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"flag"
	"./fileCache"
	"log"
	"os"
	"encoding/json"
)

type Config struct {
	API_KEY     string
	SERVER_IP   string
	SERVER_PORT string
	CACHE_TIME  int
}

var config Config

func main() {
	fmt.Println("Weather MS Loading")
	readConf()

	DEBUG := flag.Bool("debug", false, "Debug Enabled")
	flag.Parse()

	app := iris.New()
	if *DEBUG {
		app.Logger().SetLevel("debug")
		app.Use(recover.New())
		app.Use(logger.New())
	}

	app.Get("/{lat:string}/{lng:string}", func(ctx iris.Context) {
		lat := ctx.Params().Get("lat")
		lng := ctx.Params().Get("lng")
		weather, err := fileCache.GetCacheWeather(lat, lng, config.CACHE_TIME)

		if err != nil {
			weather = getWeather(lat, lng)
			fileCache.WriteWeatherToFile(lat, lng, weather)
		}

		ctx.Header("Content-Type", "application/json")
		ctx.WriteString(weather)
	})

	app.Run(iris.Addr(config.SERVER_IP+":"+config.SERVER_PORT), iris.WithoutServerError(iris.ErrServerClosed))
}

func getWeather(lat string, lng string) string {
	url := `http://api.openweathermap.org/data/2.5/forecast?lat=` + lat + `&lon=` + lng + `` + `&APPID=` + config.API_KEY + `&lang=ru&units=metric`

	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	weather := string(bs[:])
	return weather
}

func readConf() Config {
	log.Println("Reading config file")

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &config)

	return config
}
