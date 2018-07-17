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

func GetCacheWeather(lat string, lon string, cacheDuration int) (res string, err error) {
	log.Printf("Getting weather from cache. lat=%s, lon=%s", lat, lon)

	fileNameHash := GetMD5Hash(lat + lon)
	cityWeatherFilename := getCityFileName(fileNameHash)

	cacheFile, err := os.Open(cityWeatherFilename)

	if err != nil {
		return "", err
	}

	defer cacheFile.Close()

	cacheStats, err := cacheFile.Stat()

	if err != nil {
		return "", err
	}

	fileAge := time.Now().Sub(cacheStats.ModTime())
	if fileAge > time.Duration(cacheDuration)*time.Minute {
		return "", errors.New("cache expired")
	}

	cacheWeather, err := readFile(cityWeatherFilename)

	if err != nil {
		return "", err
	}

	return cacheWeather, nil
}

func WriteWeatherToFile(lat string, lon string, weather string) (err error) {
	if weather == "" {
		return errors.New("weather is empty string")
	}

	fileNameHash := GetMD5Hash(lat + lon)
	fileName := getCityFileName(fileNameHash)

	file, err := os.Create(fileName)

	if err != nil {
		log.Println("Couldn't create file " + fileName)
		log.Println("Does directory 'storage' exist?")
		return err
	}

	defer file.Close()

	file.WriteString(weather)

	return nil
}

func readFile(fileName string) (res string, err error) {
	bs, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	str := string(bs)

	return str, nil
}

func getCityFileName(city string) string {
	return defaultFileLocation + strings.ToLower(city) + ".json"
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
