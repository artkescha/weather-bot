package resolver

import (
	"context"
	"fmt"
	"log"

	"bots/telegram/weather_bot/model"
	owm "github.com/briandowns/openweathermap"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Resolver struct {
	ctx      context.Context
	token    string
	receiver tgbotapi.UpdatesChannel
}

func New(ctx context.Context, token string, receiver tgbotapi.UpdatesChannel) *Resolver {
	return &Resolver{
		ctx:      ctx,
		token:    token,
		receiver: receiver,
	}
}

func (r Resolver) Start(send func(chatID int64, messageID int, message string) error) {
	go r.run(send)
}

func (r Resolver) run(send func(chatID int64, messageID int, message string) error) {
	for {

		select {
		case update := <-r.receiver:

			if update.Message == nil {
				log.Printf("imessage is nil")
				continue
			}
			r.prepareAndSend(update, send)

		case <-r.ctx.Done():
			log.Printf("resolver context is done")
			return
		}
	}
}

func (r Resolver) prepareAndSend(update tgbotapi.Update, send func(chatID int64, messageID int, message string) error) {

	var weatherErr error

	weather, err := owm.NewForecast("5", "C", "En", r.token)
	if err != nil {
		log.Printf("new forecast failed: %s", err)
		weatherErr = fmt.Errorf("internal error: pleas try later")
	}
	err = weather.DailyByName(update.Message.Text, 0)
	if err != nil {
		log.Printf("get weather by request city name %s failed  %s", update.Message.Text, err)
		weatherErr = fmt.Errorf("internal error: get weather failed pleas try later")
	}

	weatherData, ok := weather.ForecastWeatherJson.(*owm.Forecast5WeatherData)
	if !ok {
		log.Printf("convert forecastWeatherJson to forecast5WeatherData failed %s", err)
		weatherErr = fmt.Errorf("internal error: get weather failed pleas try later")
	}

	if len(weatherData.List) == 0 {
		log.Printf("response by request city name %s is empty", update.Message.Text)
		weatherErr = fmt.Errorf(`can't find "%s" city. Try another one, for example: "Kyiv" or "Moscow"`, update.Message.Text)
	}

	myWeather, err := model.ConvertWeathersToMessage(weatherData.List)
	if err != nil {
		log.Printf("convert weather by request city name %s failed %s", update.Message.Text, err)
		weatherErr = fmt.Errorf("internal error: get weather failed pleas try later")
	}

	message := ""

	if weatherErr != nil {
		message = fmt.Sprintf("%s", weatherErr)
	} else {
		message = fmt.Sprintf("%s %s: \n%s", weatherData.City.Name, weatherData.City.Country, myWeather)
	}

	if err := send(update.Message.Chat.ID, update.Message.MessageID, message); err != nil {
		log.Printf("send message by request city name %s failed %s", update.Message.Text, err)
	}
}
