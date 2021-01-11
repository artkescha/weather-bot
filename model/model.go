package model

import (
	"fmt"
	"strings"
	"time"
)

type ForestItem struct {
	Dt          time.Time
	Temperature float64
	WindSpeed   float64
	Main        string
	Icon        string
}

func (w ForestItem) String() string {
	date := w.Dt.Format("Jan _2 15 2006")
	return fmt.Sprintf("%s \t %.1f %s %s", date, w.Temperature, w.Main, w.Icon)
}

type ForestList []ForestItem

func (list ForestList) String() string {
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

type Weather struct {
	Country string
	City    string
	List    ForestList
}
