package model

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"time"
)

func GetIcon(iconName string) (string, error) {
	switch iconName {
	case "Clear":
		return "☀", nil
	case "Rain":
		return "☔", nil
	case "Snow":
		return "❄", nil
	case "Clouds":
		return "☁", nil
	default:
		return "", fmt.Errorf("icon %s not defined", iconName)
	}
}

func ConvertWeathersToMessage(weathers []owm.Forecast5WeatherList) (WeatherList, error) {
	result := make([]Weather, 0)
	for index, weather := range weathers {
		weather, err := convertWeatherToMessage(weather)
		if err != nil {
			return []Weather{}, fmt.Errorf("weather with index %d convert failed: %s", index, err)
		}
		result = append(result, weather)
	}
	return result, nil
}

func convertWeatherToMessage(weather owm.Forecast5WeatherList) (Weather, error) {
	png, err := GetIcon(weather.Weather[0].Main)
	if err != nil {
		return Weather{}, fmt.Errorf("convert failed: %s", err)
	}
	return Weather{
		Dt:          time.Unix(int64(weather.Dt), 0),
		Temperature: weather.Main.Temp,
		WindSpeed:   weather.Wind.Speed,
		Main:        weather.Weather[0].Main,
		Icon:        png,
	}, nil
}
