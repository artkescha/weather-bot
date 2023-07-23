package resolver

import (
	"fmt"
	"github.com/artkescha/weather-bot/model"
	owm "github.com/briandowns/openweathermap"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type WeatherGetter struct {
	foreCast *owm.ForecastWeatherData
}

func NewWeatherGetter(foreCast *owm.ForecastWeatherData) WeatherGetter {
	return WeatherGetter{foreCast: foreCast}
}

func (w WeatherGetter) DailyByName(city string, days int) (model.Weather, error) {
	err := w.foreCast.DailyByName(city, days)
	if err != nil {
		return model.Weather{}, fmt.Errorf("get weather by request city name %s failed  %s", city, err)
	}
	forestData, ok := w.foreCast.ForecastWeatherJson.(*owm.Forecast5WeatherData)
	if !ok {
		return model.Weather{}, fmt.Errorf("convert forecastWeatherJson to forecast5WeatherData failed %s", err)
	}
	myWeather, err := model.ConvertForestToWeather(forestData)
	return myWeather, err
}

func (w WeatherGetter) DailyByCoordinates(location *tgbotapi.Location, days int) (model.Weather, error) {
	coordinates := &owm.Coordinates{
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
	}
	err := w.foreCast.DailyByCoordinates(coordinates, days)
	if err != nil {
		return model.Weather{}, fmt.Errorf("get weather by request coordinates city %+v failed  %s", coordinates, err)
	}
	forestData, ok := w.foreCast.ForecastWeatherJson.(*owm.Forecast5WeatherData)
	if !ok {
		return model.Weather{}, fmt.Errorf("convert forecastWeatherJson to forecast5WeatherData failed %s", err)
	}
	myWeather, err := model.ConvertForestToWeather(forestData)
	return myWeather, err
}
