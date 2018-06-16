package main

import (
	"errors"
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

// This function a string into an ID if the string is a mention.
func parseMention(str string) (string, error) {
	if len(str) < 5 || (string(str[0]) != "<" || string(str[1]) != "@" || string(str[len(str)-1]) != ">") {
		return "", errors.New("error while parsing mention, this is not an user")
	}

	res := str[2 : len(str)-1]

	// Necessary to allow nicknames.
	if string(res[0]) == "!" {
		res = res[1:]
	}

	return res, nil
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

func getRolePermissionsByName(ch *discordgo.Channel, sv *discordgo.Guild, name string) (p discordgo.PermissionOverwrite, e error) {
	//get role object for given name
	role, _ := getRoleByName(name, sv.Roles)
	return getRolePermissions(role.ID, ch.PermissionOverwrites)
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

func checkChannelsMap(data *ServerData) {
	if len(data.Channels) == 0 {
		data.Channels = make(map[string]*ChannelData)
	}
}
