package fileCache

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"log"
)

const defaultFileLocation = "storage/"

var cachedTime map[string]time.Time
var cachedWeather map[string]string

func GetCachedWeather(lat string, lon string, cacheDuration int) (res string, err error) {
	log.Printf("Getting weather from cache. lat=%s, lon=%s", lat, lon)

	latLonHash := GetMD5Hash(lat + lon)
	filename := getFileName(latLonHash)

	// Checking if cache duration expired
	fileTime, exists := cachedTime[filename]
	if exists && isExpired(fileTime, cacheDuration) {
		return "", errors.New("cache expired")
	}

	// Checking if we have requested weather stored in a variable
	variableWeather, exists := cachedWeather[filename]
	if exists {
		return variableWeather, nil
	}

	// Opening file
	cacheFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer cacheFile.Close()

	// Getting file statistics
	cacheStats, err := cacheFile.Stat()
	if err != nil {
		return "", err
	}

	// Checking if cache duration expired
	if isExpired(cacheStats.ModTime(), cacheDuration) {
		return "", errors.New("cache expired")
	}

	// Reading file content
	weather, err := readFile(filename)
	if err != nil {
		return "", err
	}

	// Saving cache file to variable
	setCachedTime(filename)
	setCachedWeather(filename, weather)

	return weather, nil
}

// Compares current time and cacheTime.
// Returns true if difference is bigger than duration.
func isExpired(cacheTime time.Time, duration int) bool {
	diff := time.Now().Sub(cacheTime)
	if diff > time.Duration(duration)*time.Minute {
		return true
	}

	return false
}

func SaveToCache(lat string, lon string, weather string) (err error) {
	latLonHash := GetMD5Hash(lat + lon)
	fileName := getFileName(latLonHash)

	file, err := os.Create(fileName)

	if err != nil {
		log.Println("Couldn't create file " + fileName)
		log.Println("Does directory 'storage' exist?")
		return err
	}

	defer file.Close()

	file.WriteString(weather)

	// Saving to variables
	setCachedTime(fileName)
	setCachedWeather(fileName, weather)

	return nil
}

func setCachedWeather(fileName string, weather string) {
	if cachedWeather == nil {
		cachedWeather = make(map[string]string)
	}
	cachedWeather[fileName] = weather
}

func setCachedTime(fileName string) {
	if cachedTime == nil {
		cachedTime = make(map[string]time.Time)
	}
	cachedTime[fileName] = time.Now()
}

func readFile(fileName string) (res string, err error) {
	bs, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	str := string(bs)

	return str, nil
}

func getFileName(city string) string {
	return defaultFileLocation + strings.ToLower(city) + ".json"
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
