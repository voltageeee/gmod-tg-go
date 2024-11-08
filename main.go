package main

import (
	"log"
	"net/http"
	"slices"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorcon/rcon"
)

// configuration of some kind
var Owner_IDs = []int64{0} // now supports multiple chat ids, add it like this: {chatid, chatid, chatid}
var Server_Addr_Port = ""
var Rcon_Pass = ""
var BotToken = ""

func sendmsg(bot *tgbotapi.BotAPI, text string) {
	for _, ID := range Owner_IDs {
		msg := tgbotapi.NewMessage(ID, text)
		msg.ParseMode = "MarkdownV2"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %d: %v", ID, err)
		}
	}
}

func handlelog(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		return
	}

	islocaladdr := strings.Contains(req.RemoteAddr, "127.0.0.1")

	if !islocaladdr {
		return
	}

	bot, err := tgbotapi.NewBotAPI(BotToken)

	if err != nil {
		log.Fatal(err)
	}

	switch req.PostFormValue("typ") {
	case "ChatMsg":
		txt := req.PostFormValue("msg")
		ply := req.PostFormValue("plr")

		sendmsg(bot, "`"+ply+": "+txt+"`")
	case "PlayerDeath":
		victim := req.PostFormValue("victim")
		attacker := req.PostFormValue("attacker")
		inflictor := req.PostFormValue("inflictor")

		sendmsg(bot, "`Player "+attacker+" killed "+victim+" using "+inflictor+"`")
	case "ConCmd":
		command := req.PostFormValue("cmd")
		player := req.PostFormValue("plr")
		arguments := req.PostFormValue("args")

		sendmsg(bot, "`Player "+player+" ran command "+command+" with args "+arguments+"`")
	case "PlayerDisconnect":
		player := req.PostFormValue("plr")
		reason := req.PostFormValue("reason")

		sendmsg(bot, "`Player "+player+" disconnected with a reason: "+reason+"`")
	case "PlayerConnect":
		player := req.PostFormValue("plr")
		ipaddr := req.PostFormValue("ipaddr")

		sendmsg(bot, "`Player "+player+" connected. IP: "+ipaddr+"`")
	}

}

func runhttp() {
	http.HandleFunc("/handlelog", handlelog)
	log.Fatal(http.ListenAndServe("127.0.0.1:5555", nil))
}

func runtgbot() {
	bot, err := tgbotapi.NewBotAPI(BotToken)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized as: %s", bot.Self.UserName)

	srv, err := rcon.Dial(Server_Addr_Port, Rcon_Pass)

	if err != nil {
		log.Print("FATAL RCON ERROR OCCURED: ", err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() || !slices.Contains(Owner_IDs, update.Message.Chat.ID) {
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
				cmd = "Fatal error occured while trying to run the command: " + err.Error()
				log.Print(err)
			}

			msg.Text = "Response from the server:\n`" + cmd + "`"
			msg.ParseMode = "MarkdownV2"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	go runtgbot()
	go runhttp()
	select {}
}
