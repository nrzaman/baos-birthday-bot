package main

import (
	_ "encoding/json"
	"flag"
	_ "flag"
	"fmt"
	_ "fmt"
	_ "github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/discord"
	"github.com/nrzaman/baos-birthday-bot/util"
	_ "io/ioutil"
	_ "net/http"
	_ "os"
	_ "os/signal"
	_ "strings"
	_ "syscall"
	_ "time"
)

// Token Variables used for command line parameters
var (
	Token string
)

// init Initializes the bot
func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

// main The main function, that involves connecting to discord and setting up message events.
func main() {
	// Extract and store birthdays from the JSON config file
	util.ExtractBirthdays()

	// Create a new Discord session using the provided bot token.
	bot, err := discord.Connect(Token)
	util.Check(err)

	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(discord.MessageCreate)
	go discord.Worker(bot)

	// Wait here until CTRL-C or other term signal is received.
	discord.WaitUntilTermination()
	err = bot.Close()
	util.Check(err)
	fmt.Println("Bot terminated.")
}
