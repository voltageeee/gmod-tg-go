package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorcon/rcon"
)

// configuration of some kind
var Owner_ID int64 = 0
var Server_Addr_Port = "0.0.0.0:27015"
var Rcon_Pass = "rconpass"
var BotToken = "bottoken"

func main() {
	bot, err := tgbotapi.NewBotAPI(BotToken)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized as: %s", bot.Self.UserName)

	srv, err := rcon.Dial(Server_Addr_Port, Rcon_Pass)

	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() || update.Message.Chat.ID != Owner_ID {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		default:
			msg.Text = "Unknown command"
		case "sendcmd":
			s := update.Message.Text
			s2, _ := strings.CutPrefix(s, "/sendcmd ")

			cmd, err := srv.Execute(s2)
			if err != nil {
				msg.Text = "A fatal error occured."
			}

			msg.Text = "Response from the server:\n`" + cmd + "`"
			msg.ParseMode = "MarkdownV2"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}
