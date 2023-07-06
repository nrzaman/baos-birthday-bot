package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/util"
	"os"
	"os/signal"
	"syscall"
)

// Connect This function starts a Discord session.
func Connect(discordToken string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + discordToken)
	util.Check(err)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent
	// Open a websocket connection to Discord and begin listening.
	err = session.Open()
	return session, err
}

// WaitUntilTermination This function listens for the user to use Ctrl+C to terminate the session.
func WaitUntilTermination() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
