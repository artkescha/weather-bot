package resolver

import (
	"context"
	"fmt"
	"log"

	"bots/telegram/weather_bot/model"
	owm "github.com/briandowns/openweathermap"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Getter interface {
	DailyByName(city string, days int) (model.Weather, error)
	DailyByCoordinates(location *owm.Coordinates, days int) (model.Weather, error)
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
				log.Printf("imessage is nil")
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

	forest, err := r.weatherGetter.DailyByName(update.Message.Text, 0)
	if err != nil {
		log.Printf("get weather failed: %s", err)
		message = "internal error: get weather failed pleas try later"
	} else {
		if forest.Country == "" {
			message = fmt.Sprintf(`can't find "%s" city. Try another one, for example: "Kyiv" or "Moscow"`, update.Message.Text)
		} else {
			message = fmt.Sprintf("%s %s: \n%s", forest.City, forest.Country, forest.List)
		}
	}
	if err := send(update.Message.Chat.ID, update.Message.MessageID, message); err != nil {
		log.Printf("send message by request city name %s failed %s", update.Message.Text, err)
	}
}
