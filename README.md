# gmod-tg-go
A simple telegram bot for interacting with your gmod (or any game that supports source rcon) server through rcon.

Make sure to update all the configuration in main.go file:
Owner_ID is your telegram user ID.
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
As for now, you can only use it to send commands to your server console using "/sendcmd" command.
For example: /sendcmd status will return smth like this:
```
hostname: Garry's Mod
version : 2024.10.29/24 9488 secure
udp/ip  : 0.0.0.0:27015  (public ip: 0.0.0.0)
steamid : [0.0.0.0] (0.0.0.0)
map     : gm_construct at: 0 x, 0 y, 0 z
players : 0 humans, 0 bots (128 max)
# userid name                uniqueid            connected ping loss state  adr
```
