package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"errors"
)


// This function parse a discord.MessageCreate into a MessageData struct.
func parseCommand(m *discordgo.MessageCreate) MessageData {
	split := strings.Fields(m.Content)

	return MessageData{m.Content[:1], strings.ToLower(split[0][1:]), split[1:], m.ID,m.ChannelID, m.Mentions,m.Author }
}

// This function a string into an ID if the string is a mention.
func parseMention(str string) (string, error) {
	if len(str) < 5 || (string(str[0]) != "<" || string(str[1]) != "@" || string(str[len(str) - 1]) != ">") {
		return "", errors.New("This is not an user.")
	}

	res := str[2:len(str)-1]

	// Necesarry to allow nicknames.
	if string(res[0]) == "!" {
		res =  res[1:len(res)-1]
	}

	return res, nil
}

// Creates a command in the given server given a name and a message.
func createCommand(data *ServerData, commandName, message string) {
	data.Commands[commandName] = &CommandData{commandName, message}
	writeServerData()
}

// Returns the ServerData of a server, given a message object.
func getServerData(s *discordgo.Session, channelId string) *ServerData {
	channel, _ := s.Channel(channelId)

	servId := channel.GuildID

	if len(Servers) == 0 {
		Servers = make(map[string]*(ServerData))
	}

	if serv, ok := Servers[servId]; ok {
		return serv
	} else {
		Servers[servId] = &ServerData{Id:servId, Key:"!"}
		return Servers[servId]
	}
}

// Checks whether a user id (String) is in a slice of users.
func userInSlice(a string, list []*discordgo.User) bool {
	for _, b := range list {
		if b.ID == a {
			return true
		}
	}
	return false
}
