package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)


func delCommander(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commanders) == 0 {
		serverData.Commanders =  make(map[string]bool)
	}

	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `" + serverData.Key + "delCommander <@User> `.")
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

func addCommander(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commanders) == 0 {
		serverData.Commanders =  make(map[string]bool)
	}

	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `" + serverData.Key + "addCommander <@User> `.")
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

func addAlbum(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Channels) == 0 {
		serverData.Channels =  make(map[string]*(ChannelData))
	}

	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `" + serverData.Key + "addalbum <AlbumID> `.")
		return
	}

	channel, ok := serverData.Channels[msg.ChannelID]

	if !ok {
		serverData.Channels[msg.ChannelID] = &ChannelData{Id: msg.ChannelID}
		channel = serverData.Channels[msg.ChannelID]
	}


	channel.Albums = append(channel.Albums, msg.Content[0])
	writeServerData()
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Added album **" + msg.Content[0] + "** to album list.")

}

func forceGetAlbum(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Channels) == 0 {
		serverData.Channels =  make(map[string]*(ChannelData))
	}

	data, _, err := img_client.GetAlbumInfo(msg.Content[0])
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, "Something went wrong while trying to retrieve an image, maybe the Imgur API is down or is there a link which is not an album?")
		return
	}
	AlbumCache[msg.Content[0]] = data
}


// Helper function to actually remove t he albums from the list.
func deleteAlbums(albumList []string, str string) []string {
	list := albumList
	for i :=0; i < len(list); i++ {
		if list[i] == str {
			list = append(list[:i], list[i+1:]...)
			i--
		}
	}

	return list
}

func delAlbum(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands =  make(map[string]*(CommandData))
	}

	if channel, ok := serverData.Channels[msg.ChannelID]; ok && len(channel.Albums) > 0 {
		serverData.Channels[msg.ChannelID].Albums = deleteAlbums(serverData.Channels[msg.ChannelID].Albums, msg.Content[0])
		writeServerData()
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Removed **" + msg.Content[0] + "** from the list of albums.")
	}
}

func addCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands =  make(map[string]*(CommandData))
	}

	if len(msg.Content) > 1 {
		createCommand(serverData, msg.Content[0], strings.Join(msg.Content[1:], " "))
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Added **"+msg.Content[0]+"** to the list of commands.")
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Adding an empty Command is quite useless, don't you think?")
	}
}

func delCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands =  make(map[string]*(CommandData))
	}

	if _, ok := serverData.Commands[msg.Content[0]]; ok {
		delete(serverData.Commands, msg.Content[0])
		writeServerData()
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Removed **" + msg.Content[0] + "** from the list of commands.")

	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "The Command **" + msg.Content[0] + "** has not been found.")
	}

}

func setKey(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Commands) == 0 {
		serverData.Commands =  make(map[string]*(CommandData))
	}

	if len(msg.Content[0]) == 1 {
		serverData.Key = msg.Content[0]
		writeServerData()
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Updated bot Key to **" + msg.Content[0] + "**.")
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Key should have a length of **1**")
	}

}

func delMe_irl(s *discordgo.Session, msg MessageData, serverData *ServerData) {
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

	if me_irl, ok := serverData.Me_irlCommands[id]; ok {
		user := me_irl.Nickname + "_irl"
		delete(serverData.Me_irlCommands, id)
		if _, ok := serverData.Commands[user]; ok {
			delete(serverData.Commands, user)
		}
	}
	writeServerData()


	_, _ = s.ChannelMessageSend(msg.ChannelID, "Removed " + msg.Key+ "me_irl Command.")
}

func addMe_irl(s *discordgo.Session, msg MessageData, serverData *ServerData) {
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

	if len(serverData.Me_irlCommands) == 0 {
		serverData.Me_irlCommands =  make(map[string]*(Me_irlData))
	}

	serverData.Me_irlCommands[id] = &Me_irlData{id, nick, content}
	writeServerData()

	_, _ = s.ChannelMessageSend(msg.ChannelID, "Added " + msg.Key+ command + ".")
}