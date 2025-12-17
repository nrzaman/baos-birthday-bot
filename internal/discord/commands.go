package bot

import "github.com/bwmarrin/discordgo"

// SlashCommands defines all the slash commands for the bot
var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "month",
		Description: "List all birthdays in the current month",
	},
	{
		Name:        "all",
		Description: "List all birthdays",
	},
	{
		Name:        "next",
		Description: "Show the next upcoming birthday",
	},
}

// RegisterCommands registers slash commands with Discord
func RegisterCommands(session *discordgo.Session, guildID string) error {
	for _, cmd := range SlashCommands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterGlobalCommands registers slash commands globally (works in all servers)
func RegisterGlobalCommands(session *discordgo.Session) error {
	for _, cmd := range SlashCommands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// CleanupCommands removes all registered commands (useful for cleanup)
func CleanupCommands(session *discordgo.Session, guildID string) error {
	commands, err := session.ApplicationCommands(session.State.User.ID, guildID)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		err := session.ApplicationCommandDelete(session.State.User.ID, guildID, cmd.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
