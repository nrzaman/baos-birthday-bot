package main

import (
	_ "encoding/json"
	"flag"
	_ "flag"
	"fmt"
	_ "fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/birthdayUtil"
	"github.com/nrzaman/baos-birthday-bot/commands"
	_ "io/ioutil"
	_ "net/http"
	"os"
	_ "os"
	"os/signal"
	_ "os/signal"
	_ "strings"
	"syscall"
	_ "syscall"
	"time"
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
	commands.ExtractBirthdays()

	// Create a new Discord session using the provided bot token.
	bot, err := Connect(Token)
	birthdayUtil.Check(err)

	for _, person := range commands.Birthdays.People {
		// Get each person's name and date of birth
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		// Check the person's birthday against current day
		if birthdayUtil.IsBirthdayCurrentDay(int(month), day) && name != "Casey" {
			_, err := bot.ChannelMessageSend("962434955680579624", fmt.Sprintf("Today is %s's birthday! Please wish them a happy birthday!", name))
			birthdayUtil.Check(err)
		} else if birthdayUtil.IsBirthdayCurrentDay(1, 6) {
			// Special handling for Casey
			_, err := bot.ChannelMessageSend("962434955680579624", fmt.Sprintf("Today is the anniversary of the Capitol Riots. Nothing else special happened today."))
			birthdayUtil.Check(err)
		}
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(commands.MessageCreate)

	// Wait here until CTRL-C or other term signal is received.
	waitUntilTermination()
	err = bot.Close()
	birthdayUtil.Check(err)
	fmt.Println("Bot terminated.")
}

// Connect This function starts a Discord session.
func Connect(discordToken string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + discordToken)
	birthdayUtil.Check(err)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent
	// Open a websocket connection to Discord and begin listening.
	err = session.Open()
	return session, err
}

// waitUntilTermination This function listens for the user to use Ctrl+C to terminate the session.
func waitUntilTermination() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
