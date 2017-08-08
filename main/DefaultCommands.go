package main

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"

	"time"
)

// Retrieves the meIrl of a given user.
func meIrl(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if res, ok := serverData.meIrlCommands[msg.Author.ID]; ok {
		_, _ = s.ChannelMessageSend(msg.ChannelID, res.Content)
	} else {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Sorry, you do not have a meIrl. ")
	}
}

// Handles the i or image command.
// Checks whether an album is actually in cache before making an imgur API call.
func getImage(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(serverData.Channels) == 0 {
		serverData.Channels = make(map[string]*(ChannelData))
	}

	// Check there is server data for the given channel and check if there is at least 1 album.
	if channel, ok := serverData.Channels[msg.ChannelID]; ok && len(channel.Albums) > 0 {
		// Get a random index and get the Album ID on that index.
		nmbr := rand.Intn(time.Now().Nanosecond()) % len(channel.Albums)
		id := channel.Albums[nmbr]

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

func help(s *discordgo.Session, msg MessageData, serverData *ServerData) {

}
