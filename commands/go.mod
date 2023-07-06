module github.com/nrzaman/baos-birthday-bot/commands

go 1.20

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/nrzaman/baos-birthday-bot/birthdayUtil v0.0.0-20230706021010-b6458f5c8395
)

replace github.com/nrzaman/baos-birthday-bot/birthdayUtil => ../birthdayUtil