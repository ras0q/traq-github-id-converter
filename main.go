package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
		// ex: /register <traQ ID> <GitHub ID>
		if strings.HasPrefix(p.Message.PlainText, "/register") {
			args := strings.Split(p.Message.PlainText, " ")
			if len(args) != 3 {
				mustPostMessage(bot, "Usage: `/register <traQ ID> <GitHub ID>`")
			}

			traqID := args[1]
			githubID := args[2]

			fmt.Println(traqID, githubID) // TODO: register to DB

			mustPostMessage(bot, "Registered!")
		}

		mustPostMessage(bot, "@Ras")
	})

	panicOnError(bot.Start())
}

func mustPostMessage(bot *traqwsbot.Bot, content string) {
	embed := true

	_, _, err := bot.API().
		MessageApi.
		PostMessage(context.Background(), channelID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: content,
			Embed:   &embed,
		}).
		Execute()
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
