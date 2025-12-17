package channel

import (
	"encoding/json"

	"github.com/nrzaman/baos-birthday-bot/internal/interfaces"
	"github.com/nrzaman/baos-birthday-bot/util"
)

// Service handles channel-related operations with injected dependencies
type Service struct {
	fileReader interfaces.FileReader
	channels   util.Channels
}

// NewService creates a new channel Service with the given dependencies
func NewService(fileReader interfaces.FileReader) *Service {
	return &Service{
		fileReader: fileReader,
	}
}

// LoadChannels loads channels from the config file
func (s *Service) LoadChannels(configPath string) error {
	byteResult, err := s.fileReader.ReadFile(configPath)
	if err != nil {
		return err
	}

	var channels util.Channels
	err = json.Unmarshal(byteResult, &channels)
	if err != nil {
		return err
	}

	s.channels = channels
	return nil
}

// GetGeneralChannelID retrieves the general channel ID
func (s *Service) GetGeneralChannelID() string {
	for _, channel := range s.channels.Channel {
		if channel.Name == "general" {
			return channel.ID
		}
	}
	return ""
}

// GetChannelByName retrieves a channel ID by name
func (s *Service) GetChannelByName(name string) string {
	for _, channel := range s.channels.Channel {
		if channel.Name == name {
			return channel.ID
		}
	}
	return ""
}
