package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// This function parse a discord.MessageCreate into a MessageData struct.
func parseCommand(m *discordgo.MessageCreate) MessageData {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error with: `" + m.Content + "`, I should fix dis...")
		}
	}()

	split := strings.Fields(m.Content)
	Key := m.Content[:1]
	Command := strings.ToLower(split[0][1:])
	Content := split[1:]
	return MessageData{Key, Command, Content, m.ID, m.ChannelID, m.Mentions, m.Author}
}

// This function a string into an ID if the string is a mention.
func parseMention(str string) (string, error) {
	if len(str) < 5 || (string(str[0]) != "<" || string(str[1]) != "@" || string(str[len(str)-1]) != ">") {
		return "", errors.New("This is not an user.")
	}

	res := str[2:]

	// Necessary to allow nicknames.
	if string(res[0]) == "!" {
		res = res[1 : len(res)-1]
	}

	return res, nil
}

// Creates a command in the given server given a name and a message.
func createCommand(data *ServerData, commandName, message string) {
	data.Commands[commandName] = &CommandData{strings.ToLower(commandName), message}
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
		Servers[servId] = &ServerData{Id: servId, Key: "!"}
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

// Gets the specific role by name out of a role list.
func getRoleByName(name string, roles []*discordgo.Role) (r discordgo.Role, e error) {
	for _, elem := range roles {
		if elem.Name == name {
			r = *elem
			return
		}
	}
	e = errors.New("Role name not found in the specified role array: " + name)
	return
}

// Gets the permission override object from a role id.
func getRolePermissions(id string, perms []*discordgo.PermissionOverwrite) (p discordgo.PermissionOverwrite, e error) {
	for _, elem := range perms {
		if elem.ID == id {
			p = *elem
			return
		}
	}
	e = errors.New("Permissions not found in the specified role: " + id)
	return
}
