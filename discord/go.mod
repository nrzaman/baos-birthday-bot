module github.com/nrzaman/baos-birthday-bot/discord

go 1.20

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/nrzaman/baos-birthday-bot/util v0.0.0-20230706021010-b6458f5c8395
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)

replace github.com/nrzaman/baos-birthday-bot/util => ../util
