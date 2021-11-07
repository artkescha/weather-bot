package resolver

import (
	"fmt"
	"log"

	"github.com/artkescha/weather-bot/model"
	owm "github.com/briandowns/openweathermap"
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
	return myWeather, nil
}

func (w WeatherGetter) DailyByCoordinates(location *owm.Coordinates, days int) (model.Weather, error) {
	log.Print("implement me later")
	return model.Weather{}, nil
}
