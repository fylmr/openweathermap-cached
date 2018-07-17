# openweathermap-cached
Simple server that gets info from OpenWeatherMap and caches it, so you wouldn't run out of OWM requests.

With some pretty easy modifications, you can use it for any other API and simple web pages that doesn't require constant refreshing.  

## How to use

Place **config.json** in the same folder as main.go.

**config.json** should have structure like so:
````
{
  "WEATHER_HOST": "http://api.openweathermap.org/data/2.5/forecast",
  "API_KEY": "YOUR_API_KEY",
  "SERVER_IP": "127.0.0.1",
  "SERVER_PORT": "9001",
  "CACHE_TIME": 15,
  "LANGUAGE": "ru",
  "UNITS": "metric"
}
````
To read all the information about API values, please visut [OpenWeatherMap site](https://openweathermap.org/forecast5). 
