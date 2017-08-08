package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Deletes a commander from the list of commanders.
// First argument should be a mention to the person who should be deleted.
func delCommander(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commanders) == 0 {
		serverData.Commanders = make(map[string]bool)
	}

	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `"+serverData.Key+"delCommander <@User> `.")
		return
	}

	userID, err := parseMention(msg.Content[0])

	if err != nil {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Something went wrong while trying to add user as commander")
		return
	}

	serverData.Commanders[userID] = false
	writeServerData()
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Removed commander.")
}

// Add a commander to the list of commanders.
// First argument should be a mention to the person who should be added.
func addCommander(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commanders) == 0 {
		serverData.Commanders = make(map[string]bool)
	}

	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `"+serverData.Key+"addCommander <@User> `.")
		return
	}

	userID, err := parseMention(msg.Content[0])

	if err != nil {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Something went wrong while trying to add user as commander")
		return
	}

	serverData.Commanders[userID] = true
	writeServerData()
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Added commander.")

}

// Adds an album to the list of albums.
func addAlbum(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Channels) == 0 {
		serverData.Channels = make(map[string]*(ChannelData))
	}

	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `"+serverData.Key+"addalbum <AlbumID> `.")
		return
	}

	channel, ok := serverData.Channels[msg.ChannelID]

	if !ok {
		serverData.Channels[msg.ChannelID] = &ChannelData{ID: msg.ChannelID}
		channel = serverData.Channels[msg.ChannelID]
	}

	channel.Albums = append(channel.Albums, msg.Content[0])
	writeServerData()
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Added album **"+msg.Content[0]+"** to album list.")

}

// Refresh the data of the album which is cached.
func forceGetAlbum(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Channels) == 0 {
		serverData.Channels = make(map[string]*(ChannelData))
	}

	data, _, err := imgClient.GetAlbumInfo(msg.Content[0])
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, "Something went wrong while trying to retrieve an image, maybe the Imgur API is down or is there a link which is not an album?")
		return
	}
	AlbumCache[msg.Content[0]] = data
}

// Helper function to actually remove the albums from the list.
func deleteAlbums(albumList []string, str string) []string {
	list := albumList
	for i := 0; i < len(list); i++ {
		if list[i] == str {
			list = append(list[:i], list[i+1:]...)
			i--
		}
	}

	return list
}

// Remove a specific album from the list of albums.
func delAlbum(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands = make(map[string]*(CommandData))
	}

	if channel, ok := serverData.Channels[msg.ChannelID]; ok && len(channel.Albums) > 0 {
		serverData.Channels[msg.ChannelID].Albums = deleteAlbums(serverData.Channels[msg.ChannelID].Albums, msg.Content[0])
		writeServerData()
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Removed **"+msg.Content[0]+"** from the list of albums.")
	}
}

// Add a command to the list  of commands.
// First argument should be the name, the rest should be the content.
// Will overwrite existing commands!
func addCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands = make(map[string]*(CommandData))
	}

	if len(msg.Content) > 1 {
		createCommand(serverData, msg.Content[0], strings.Join(msg.Content[1:], " "))
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Added **"+msg.Content[0]+"** to the list of commands.")
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Adding an empty Command is quite useless, don't you think?")
	}
}

// Removes a command from the list of commands
// First argument should be the name of the command that should be removed
// Nothing happens if command is unknown.
func delCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands = make(map[string]*(CommandData))
	}

	if _, ok := serverData.Commands[msg.Content[0]]; ok {
		delete(serverData.Commands, msg.Content[0])
		writeServerData()
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Removed **"+msg.Content[0]+"** from the list of commands.")

	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "The Command **"+msg.Content[0]+"** has not been found.")
	}

}

// Change the command key of the server.
// Key should be the first argument and only 1 char long.
func setKey(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands = make(map[string]*(CommandData))
	}

	if len(msg.Content[0]) == 1 {
		serverData.Key = msg.Content[0]
		writeServerData()
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Updated bot Key to **"+msg.Content[0]+"**.")
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Key should have a length of **1**")
	}

}

// Remove a meIrl command from a specific person.
func delMeIrl(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	usage := "Usage: `" + msg.Key + "delme_irl <@User> `"

	if len(msg.Content) < 1 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, usage)
		return
	}

	id, err := parseMention(msg.Content[0])

	if err != nil {
		_, _ = s.ChannelMessageSend(msg.ChannelID, usage)
		return
	}

	if meIrl, ok := serverData.meIrlCommands[id]; ok {
		user := meIrl.Nickname + "_irl"
		delete(serverData.meIrlCommands, id)
		if _, ok := serverData.Commands[user]; ok {
			delete(serverData.Commands, user)
		}
	}
	writeServerData()

	_, _ = s.ChannelMessageSend(msg.ChannelID, "Removed "+msg.Key+"meIrl Command.")
}

// Add a meIrl command from a specific person
func addMeIrl(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	usage := "Usage: `" + msg.Key + "addme_irl <@User> <Nickname> <Content> `"
	if len(msg.Content) < 3 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, usage)
		return
	}

	id, err := parseMention(msg.Content[0])

	if err != nil {
		_, _ = s.ChannelMessageSend(msg.ChannelID, usage)
		return
	}

	nick := msg.Content[1]
	content := strings.Join(msg.Content[2:], " ")

	command := nick + "_irl"
	createCommand(serverData, command, content)

	if len(serverData.meIrlCommands) == 0 {
		serverData.meIrlCommands = make(map[string]*(MeIrlData))
	}

	serverData.meIrlCommands[id] = &MeIrlData{id, nick, content}
	writeServerData()

	_, _ = s.ChannelMessageSend(msg.ChannelID, "Added "+msg.Key+command+".")
}

// Lock a channel, so the @everyone role won't be able to talk in the channel.
func lockChannel(s *discordgo.Session, msg MessageData, serverData *ServerData) {

	//get channel object
	ch, err := s.Channel(msg.ChannelID)
	if err != nil {
		fmt.Println("Couldn't find channel with following id: ", msg.ChannelID)
		return
	}

	//get server object
	sv, err := s.Guild(serverData.ID)
	if err != nil {
		fmt.Println("Couldn't find channel with following id: ", serverData.ID)
		return
	}

	//get @everyone role object
	role, err := getRoleByName("@everyone", sv.Roles)
	if err != nil {
		fmt.Println("Couldn't find @everyone role of the following server: ", serverData.ID)
		return
	}

	//get @everyone permissions
	everyonePerms, err := getRolePermissions(role.ID, ch.PermissionOverwrites)
	if err != nil {
		fmt.Println("Couldn't get @everyone permissions of the following server: ", serverData.ID)
		return
	}

	//deny sending messages and update it
	err = s.ChannelPermissionSet(ch.ID, everyonePerms.ID, everyonePerms.Type, everyonePerms.Allow^0x800, everyonePerms.Deny|0x800)
	if err != nil {
		fmt.Println("Error unlocking channel: ", err)
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "This channel is now locked.")
	}
}

// Unlock the channel, so the @everyone role will be allowed to talk again.
func unlockChannel(s *discordgo.Session, msg MessageData, serverData *ServerData) {

	//get channel object
	ch, err := s.Channel(msg.ChannelID)
	if err != nil {
		fmt.Println("Couldn't find channel with following id: ", msg.ChannelID)
		return
	}

	//get server object
	sv, err := s.Guild(serverData.ID)
	if err != nil {
		fmt.Println("Couldn't find channel with following id: ", serverData.ID)
		return
	}

	//get @everyone role object
	role, err := getRoleByName("@everyone", sv.Roles)
	if err != nil {
		fmt.Println("Couldn't find @everyone role of the following server: ", serverData.ID)
		return
	}

	//get @everyone permissions
	everyonePerms, err := getRolePermissions(role.ID, ch.PermissionOverwrites)
	if err != nil {
		fmt.Println("Couldn't get @everyone permissions of the following server: ", serverData.ID)
		return
	}

	//allow sending messages and update it
	err = s.ChannelPermissionSet(ch.ID, everyonePerms.ID, everyonePerms.Type, everyonePerms.Allow|0x800, everyonePerms.Deny^0x800)
	if err != nil {
		fmt.Println("Error unlocking channel: ", err)
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "This channel is now unlocked.")
	}
}
