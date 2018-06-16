package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// This function parse a discord.MessageCreate into a SentMessageData struct.
func parseMessage(m *discordgo.MessageCreate) SentMessageData {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Error with: `" + m.Content + "`, I should fix this...")
		}
	}()

	split := strings.Fields(m.Content)
	Key := m.Content[:1]
	Command := strings.ToLower(split[0][1:])
	Content := split[1:]
	return SentMessageData{Key, Command, Content, m.ID, m.ChannelID, m.Mentions, m.Author}
}

// Returns the ServerData of a server, given a message object.
func getServerData(s *discordgo.Session, channelID string) *ServerData {
	channel, _ := s.Channel(channelID)

	servID := channel.GuildID

	if len(Servers) == 0 {
		Servers = make(map[string]*ServerData)
	}

	if serv, ok := Servers[servID]; ok {
		return serv
	}

	Servers[servID] = &ServerData{ID: servID, Key: "!"}
	return Servers[servID]

}

// Creates a command in the given server given a name and a message.
func createCommand(data *ServerData, commandName, message string) {
	data.CustomCommands[commandName] = &CommandData{strings.ToLower(commandName), message}
	writeServerData()
}

func checkCommandsMap(data *ServerData) {
	if len(data.CustomCommands) == 0 {
		data.CustomCommands = make(map[string]*CommandData)
	}
}
