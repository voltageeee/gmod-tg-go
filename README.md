# gmod-tg-go
A simple telegram bot for interacting with your gmod (or any game that supports source rcon) server through rcon.

Make sure to update all the configuration in main.go file:
```
var Owner_ID int64 = 0
var Server_Addr_Port = "0.0.0.0:27015"
var Rcon_Pass = "rconpass"
var BotToken = "bottoken"
```
And to install all the dependencies:
```
go get github.com/gorcon/rcon
go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5
```
