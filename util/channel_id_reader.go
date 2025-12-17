package util

import (
	"encoding/json"
)

var DiscordChannels Channels
var GeneralChannelID string

// ExtractChannels This function extracts channels from a config file to obtain channel IDs.
func ExtractChannels() {
	// Read all contents
	var byteResult = Extract("./config/channels.json")

	// Create result variable
	var channels Channels

	// Store contents
	var err = json.Unmarshal(byteResult, &channels)
	Check(err)
	DiscordChannels = channels
	GetGeneralChannelID()
}

// GetGeneralChannelID Retrieves and stores the general channel ID.
func GetGeneralChannelID() {
	for _, channel := range DiscordChannels.Channel {
		name := channel.Name
		ID := channel.ID

		if name == "general" {
			GeneralChannelID = ID
		}
	}
}
