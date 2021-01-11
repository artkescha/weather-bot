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

func ConvertForestToWeather(forest *owm.Forecast5WeatherData) (Weather, error) {
	result := make([]ForestItem, 0)
	for index, weather := range forest.List {
		weather, err := convertForestDataToItem(weather)
		if err != nil {
			return Weather{}, fmt.Errorf("forest item with index %d convert failed: %s", index, err)
		}
		result = append(result, weather)
	}
	return Weather{
		Country: forest.City.Country,
		City:    forest.City.Name,
		List:    result}, nil
}

func convertForestDataToItem(weather owm.Forecast5WeatherList) (ForestItem, error) {
	png, err := GetIcon(weather.Weather[0].Main)
	if err != nil {
		return ForestItem{}, fmt.Errorf("get icon for data %s failed: %s", weather.Weather[0].Main, err)
	}
	return ForestItem{
		Dt:          time.Unix(int64(weather.Dt), 0),
		Temperature: weather.Main.Temp,
		WindSpeed:   weather.Wind.Speed,
		Main:        weather.Weather[0].Main,
		Icon:        png,
	}, nil
}
