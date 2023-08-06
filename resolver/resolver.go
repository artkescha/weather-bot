package resolver

import (
	"context"
	"fmt"
	"github.com/artkescha/weather-bot/model"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Getter interface {
	DailyByName(city string, days int) (model.Weather, error)
	DailyByCoordinates(location *tgbotapi.Location, days int) (model.Weather, error)
}

type Resolver struct {
	ctx           context.Context
	token         string
	weatherGetter Getter
	receiver      tgbotapi.UpdatesChannel
}

func New(ctx context.Context, token string, weatherGetter Getter, receiver tgbotapi.UpdatesChannel) *Resolver {
	return &Resolver{
		ctx:           ctx,
		token:         token,
		weatherGetter: weatherGetter,
		receiver:      receiver,
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
				log.Printf("message is nil")
				continue
			}
			r.getForestAndSend(update, send)
		case <-r.ctx.Done():
			log.Printf("resolver context is done")
			return
		}
	}
}

func (r Resolver) getForestAndSend(update tgbotapi.Update, send func(chatID int64, messageID int, message string) error) {
	message := ""
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			if err := send(update.Message.Chat.ID, update.Message.MessageID, "click button 'My city' or send city name"); err != nil {
				log.Printf("send message by request city name %s failed %s", update.Message.Text, err)
			}
			return
		}
	}
	forest, err := r.weatherForecast(update.Message, 0)
	if err != nil {
		log.Printf("get weather failed: %s", err)
		message = fmt.Sprintf("get weather forecast failed, reason %s pleas try again later", err)
	} else {
		if len(forest.List) == 0 {
			message = fmt.Sprintf(`can't find "%s" city. Try another one, for example: "Kyiv" or "Moscow"`, update.Message.Text)
		} else {
			message = fmt.Sprintf("%s %s: \n%s", forest.City, forest.Country, forest.List)
		}
	}
	if err := send(update.Message.Chat.ID, update.Message.MessageID, message); err != nil {
		log.Printf("send message by request city name %s failed %s", update.Message.Text, err)
	}
}

func (r Resolver) weatherForecast(message *tgbotapi.Message, days int) (model.Weather, error) {
	switch {
	case message.Location != nil:
		return r.weatherGetter.DailyByCoordinates(message.Location, days)
	default:
		return r.weatherGetter.DailyByName(message.Text, days)
	}
}
