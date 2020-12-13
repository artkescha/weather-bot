package model

import (
	"fmt"
	"strings"
	"time"
)

type Weather struct {
	Dt          time.Time
	Temperature float64
	WindSpeed   float64
	Main        string
	Icon        string
}

func (w Weather) String() string {
	date := w.Dt.Format("Jan _2 15 2006")
	return fmt.Sprintf("%s \t %.1f %s %s", date, w.Temperature, w.Main, w.Icon)
}

type WeatherList []Weather

func (list WeatherList) String() string {
	var result strings.Builder
	prevPoint := 0
	for _, weather := range list {
		currentPoint := weather.Dt.Day() + int(weather.Dt.Month()) + weather.Dt.Year()
		if currentPoint > prevPoint {
			result.WriteString(fmt.Sprintf("\n%d %v\n", weather.Dt.Day(), weather.Dt.Month()))
		}
		result.WriteString(fmt.Sprintf("%s\n", weather))

		prevPoint = currentPoint

	}
	return result.String()
}
