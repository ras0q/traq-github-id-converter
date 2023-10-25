package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

var (
	channelID   = os.Getenv("TRAQ_CHANNEL_ID")
	accessToken = os.Getenv("TRAQ_ACCESS_TOKEN")
	mysqlConfig = mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: accessToken,
	})
	panicOnError(err)

	bot.OnError(func(msg string) {
		log.Println("Received Error:", msg)
	})

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	panicOnError(err)

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		// ex: /register <traQ ID> <GitHub ID>
		if strings.HasPrefix(p.Message.PlainText, "/register") {
			args := strings.Split(p.Message.PlainText, " ")
			if len(args) != 3 {
				mustPostMessage(bot, "Usage: `/register <traQ ID> <GitHub ID>`")
			}

			traqID := args[1]
			githubID := args[2]

			_, err := db.ExecContext(
				context.Background(),
				"INSERT INTO `users` (`traq_id`, `github_id`) VALUES (?, ?)",
				traqID, githubID,
			)
			if err != nil {
				mustPostMessage(bot, "Failed to register")
				return
			}

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
