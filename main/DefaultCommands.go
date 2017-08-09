package main

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"

	"fmt"
	"reflect"
	"strings"
	"time"
	"strconv"
)

// Retrieves the meIrlCommand of a given user.
func meIrlCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if res, ok := serverData.MeIrlCommands[msg.Author.ID]; ok {
		_, _ = s.ChannelMessageSend(msg.ChannelID, res.Content)
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Sorry, you do not have a me_irl. ")
	}
}

// Handles the i or image command.
// Checks whether an album is actually in cache before making an imgur API call.
func getImageCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Channels) == 0 {
		serverData.Channels = make(map[string]*(ChannelData))
	}

	// Check there is server data for the given channel and check if there is at least 1 album.
	if channel, ok := serverData.Channels[msg.ChannelID]; ok && len(channel.Albums) > 0 {
		// Get a random index and get the Album ID on that index.
		randomImageID := rand.Intn(time.Now().Nanosecond()) % len(channel.Albums)
		id := channel.Albums[randomImageID]

		// Get the data from the cache.
		data, ok := AlbumCache[id]

		// If album is not already in cache, retrieve it from the Imgur servers.
		if !ok {
			var err error
			data, _, err = imgClient.GetAlbumInfo(id)
			if err != nil {
				s.ChannelMessageSend(msg.ChannelID, "Something went wrong while trying to retrieve an image, maybe the Imgur API is down or **"+id+"** is not an album?")
				return
			}
			AlbumCache[id] = data

		}

		// Get a random image from the album and get the link of said image.
		rndImg := rand.Intn(time.Now().Nanosecond()) % len(data.Images)
		linkToImg := data.Images[rndImg].Link

		s.ChannelMessageSend(msg.ChannelID, linkToImg)
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "This channel does not have any albums, add an album using `"+serverData.Key+"addalbum <AlbumID> `.")
	}

}

func rollCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	var maxRoll int64 = 100

	if len(msg.Content) > 0 {
		newMaxRoll, err := strconv.ParseInt(msg.Content[0], 10, 64)
		if err == nil && newMaxRoll > 0{
			maxRoll = newMaxRoll
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	_, _ = s.ChannelMessageSend(msg.ChannelID, strconv.FormatInt(rand.Int63n(maxRoll + 1), 10))

}

// Send a private message with the basic info of the bot
func helpCommand(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	channel, err := s.UserChannelCreate(msg.Author.ID)
	if err != nil {
		fmt.Println("Unable to open private channel with user, ", msg.Author.ID)
	}

	s.ChannelMessageSend(channel.ID,
		"Heya, This is OwOBot written by GenDoNL. \n"+
			"For the list of commands use **"+serverData.Key+"commandlist** in the server. \n\n"+
			"The source code of the bot can be found here: https://github.com/GenDoNL/GoDiscordBot")
}

// Get all the message
func commandListCommands(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	channel, err := s.UserChannelCreate(msg.Author.ID)
	if err != nil {
		fmt.Println("Unable to open private channel with user, ", msg.Author.ID)
	}

	keys := reflect.ValueOf(serverData.Commands).MapKeys()
	strkeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}

	message := strings.Join(strkeys, ", ")

	if len(message) > 1950 {
		message = message[0:1950] + "...truncated"
	}

	s.ChannelMessageSend(channel.ID,
		"```"+message+" ```")
}
