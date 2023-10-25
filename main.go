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

	_, err = db.ExecContext(
		context.Background(),
		"CREATE TABLE IF NOT EXISTS `users` (`traq_id` VARCHAR(36) NOT NULL, `github_id` VARCHAR(39) NOT NULL, PRIMARY KEY (`traq_id`))",
	)
	panicOnError(err)

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		ctx := context.Background()

		// ex: /register <traQ ID> <GitHub ID>
		if strings.HasPrefix(p.Message.PlainText, "/register") {
			args := strings.Split(p.Message.PlainText, " ")
			if len(args) != 3 {
				mustPostMessage(ctx, bot, "Usage: `/register <traQ ID> <GitHub ID>`", p.Message.ChannelID)
			}

			traqID := args[1]
			githubID := args[2]

			_, err := db.ExecContext(
				ctx,
				"INSERT INTO `users` (`traq_id`, `github_id`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `github_id` = ?",
				traqID, githubID,
			)
			if err != nil {
				mustPostMessage(ctx, bot, "Failed to register", p.Message.ChannelID)
				return
			}

			mustPostMessage(ctx, bot, "Registered!", p.Message.ChannelID)
		}

		mustPostMessage(ctx, bot, "@Ras", p.Message.ChannelID)
	})

	panicOnError(bot.Start())
}

func mustPostMessage(ctx context.Context, bot *traqwsbot.Bot, content string, channelID string) {
	embed := true

	_, _, err := bot.API().
		MessageApi.
		PostMessage(ctx, channelID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: content,
			Embed:   &embed,
		}).
		Execute()
	if err != nil {
		log.Println("Failed to post message:", err)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
