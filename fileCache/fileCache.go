package fileCache

import (
	"os"
	"time"
	"io/ioutil"
	"strings"
	"errors"
	"crypto/md5"
	"encoding/hex"
)

const defaultFileLocation = "storage/"

func GetCacheWeather(lat string, lng string, cacheDuration int) (res string, err error) {
	fileNameHash := GetMD5Hash(lat + lng)
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

func WriteWeatherToFile(lat string, lng string, weather string) (err error) {
	fileNameHash := GetMD5Hash(lat + lng)
	fileName := getCityFileName(fileNameHash)
	file, err := os.Create(fileName)

	if err != nil {
		println()
		println("Cant create file: " + fileName)
		println("Does directory 'storage' exists?")
		println()
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
