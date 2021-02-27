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
	prevDate := time.Time{}
	for _, weather := range list {

		currentDate := time.Date(weather.Dt.Year(), weather.Dt.Month(), weather.Dt.Day(), 0, 0, 0, 0, weather.Dt.Location())

		if currentDate.Sub(prevDate) >= (24 * time.Hour) {
			result.WriteString(fmt.Sprintf("\n%d %v\n", weather.Dt.Day(), weather.Dt.Month()))
		}

		result.WriteString(fmt.Sprintf("%s\n", weather))
		prevDate = currentDate
	}
	return result.String()
}

type Weather struct {
	Country string
	City    string
	List    ForestList
}
