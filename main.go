package main

import (
	"context"
	"log"
	"os"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

var (
	channelID   = os.Getenv("TRAQ_CHANNEL_ID")
	accessToken = os.Getenv("TRAQ_ACCESS_TOKEN")
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: accessToken,
	})
	panicOnError(err)

	bot.OnError(func(msg string) {
		log.Println("Received Error:", msg)
	})

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		embed := true
		_, _, err := bot.API().
			MessageApi.
			PostMessage(context.Background(), channelID).
			PostMessageRequest(traq.PostMessageRequest{
				Content: "@Ras",
				Embed:   &embed,
			}).
			Execute()
		if err != nil {
			log.Println("Failed to post message:", err)
		}
	})
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
