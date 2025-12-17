package interfaces

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// TimeProvider provides time-related functionality that can be mocked in tests
type TimeProvider interface {
	Now() time.Time
	Month() time.Month
	Day() int
}

// FileReader provides file reading functionality that can be mocked in tests
type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

// DiscordClient provides Discord-related functionality that can be mocked in tests
type DiscordClient interface {
	SendMessage(channelID string, message string) error
	AddHandler(handler interface{})
	Close() error
}

// DiscordSession wraps the real discordgo.Session to implement DiscordClient
type DiscordSession struct {
	Session *discordgo.Session
}

func (ds *DiscordSession) SendMessage(channelID string, message string) error {
	_, err := ds.Session.ChannelMessageSend(channelID, message)
	return err
}

func (ds *DiscordSession) AddHandler(handler interface{}) {
	ds.Session.AddHandler(handler)
}

func (ds *DiscordSession) Close() error {
	return ds.Session.Close()
}
