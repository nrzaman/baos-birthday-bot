package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/nrzaman/baos-birthday-bot/internal/birthday"
	"github.com/nrzaman/baos-birthday-bot/internal/database"
	bot "github.com/nrzaman/baos-birthday-bot/internal/discord"
	"github.com/nrzaman/baos-birthday-bot/internal/interfaces"
	"github.com/nrzaman/baos-birthday-bot/internal/providers"
)

// main The main function with dependency injection
func main() {
	// Load configuration from environment variables
	discordToken := os.Getenv("DISCORD_BIRTHDAY_BOT_TOKEN")
	if discordToken == "" {
		log.Fatal("DISCORD_BIRTHDAY_BOT_TOKEN environment variable is required")
	}

	generalChannelID := os.Getenv("DISCORD_BIRTHDAY_CHANNEL_ID")
	if generalChannelID == "" {
		log.Fatal("DISCORD_BIRTHDAY_CHANNEL_ID environment variable is required")
	}

	// Create real implementations of our dependencies
	timeProvider := &providers.RealTimeProvider{}

	// Open database
	db, err := database.New("./birthdays.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Create services with injected dependencies
	birthdayService := birthday.NewServiceDB(timeProvider, db)

	// Create Discord session
	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsMessageContent

	// Wrap Discord session in our interface
	discordClient := &interfaces.DiscordSession{Session: session}

	// Create handler with dependencies
	handler := bot.NewHandler(discordClient, birthdayService)

	// Register slash command handler
	session.AddHandler(handler.HandleSlashCommand)

	// Open websocket connection
	if err := session.Open(); err != nil {
		log.Fatalf("Failed to open connection: %v", err)
	}

	// Register slash commands globally (works in all servers)
	fmt.Println("Registering slash commands...")
	if err := bot.RegisterGlobalCommands(session); err != nil {
		log.Printf("Warning: Failed to register slash commands: %v", err)
		log.Println("Slash commands may not work, but legacy !commands will still work")
	} else {
		fmt.Println("Slash commands registered successfully!")
		fmt.Println("Available commands: /month, /all, /next")
	}

	// Start worker in background
	worker := bot.NewWorker(discordClient, birthdayService, timeProvider, generalChannelID)
	go worker.Start()

	// Wait for termination signal
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanup
	fmt.Println("Shutting down...")
	worker.Stop()

	// Note: We don't clean up slash commands here as they persist across bot restarts
	// To remove them manually, you can call bot.CleanupCommands() if needed

	if err := session.Close(); err != nil {
		log.Printf("Error closing session: %v", err)
	}
	fmt.Println("Bot terminated.")
}
