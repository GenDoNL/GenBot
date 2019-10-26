package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/saucenao-go"
	"github.com/koffeinsource/go-imgur"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ImageModule struct {
	ImageCommands map[string]Command
}

func (cmd *ImageModule) setup() {
	imgClient = new(imgur.Client)
	imgClient.HTTPClient = new(http.Client)
	imgClient.Log = log
	imgClient.ImgurClientID = BotConfig.ImgurToken

	AlbumCache = make(map[string]*imgur.AlbumInfo)

	cmd.ImageCommands = map[string]Command{}

	addAlbumCommand := Command{
		Name:        "addalbum",
		Description: "This command adds an album to this channel. Images can be retrieved at random using the `i` or `image` command.",
		Usage:       "Usage: `%saddalbum <imgur album id>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     addAlbumCommand,
	}
	cmd.ImageCommands["addalbum"] = addAlbumCommand

	delAlbumCommand := Command{
		Name:        "delalbum",
		Description: "This command removes an album from this channel.",
		Usage:       "Usage: `%sdelalbum <imgur album id>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     delAlbumCommand,
	}
	cmd.ImageCommands["delalbum"] = delAlbumCommand

	imageCommand := Command{
		Name:        "image",
		Description: "This command sends a random image from an album added by `addalbum`",
		Usage:       "Usage: `%si`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     getImageCommand,
	}
	cmd.ImageCommands["i"] = imageCommand
	cmd.ImageCommands["image"] = imageCommand

	sauceCommand := Command{
		Name:        "sauce",
		Description: "Provides the source of an image. A direct link to an image should be provided.",
		Usage:       "Usage: `%ssource <URL>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     getSauceCommand,
	}
	cmd.ImageCommands["sauce"] = sauceCommand
	cmd.ImageCommands["source"] = sauceCommand

}

func (cmd *ImageModule) retrieveCommands() map[string]Command {
	return cmd.ImageCommands
}

func (cmd *ImageModule) retrieveHelp() string {
	return ""
}

func (cmd *ImageModule) execute(s *discordgo.Session, m *discordgo.MessageCreate, msg SentMessageData, serverData *ServerData) {
	if serverData.Key != msg.Key {
		return
	}

	if command, ok := cmd.ImageCommands[msg.Command]; ok {
		isCommander, ok := serverData.Commanders[m.Author.ID]
		perm, _ := s.UserChannelPermissions(msg.Author.ID, msg.ChannelID)

		// Check if user has the correct permission or whether user is a commander
		if perm&command.Permission == command.Permission || (ok && isCommander) {
			log.Infof("Executing command `%s` in server `%s` ", command.Name, serverData.ID)
			command.Execute(command, s, msg, serverData)
		} else {
			log.Infof("Use of %s command denied for permission level %d", command.Name, perm)
		}
	}
}

// Adds an album to the list of albums.
func addAlbumCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	checkChannelsMap(data)

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	channel, ok := data.Channels[msg.ChannelID]

	if !ok {
		data.Channels[msg.ChannelID] = &ChannelData{ID: msg.ChannelID}
		channel = data.Channels[msg.ChannelID]
	}

	channel.Albums = append(channel.Albums, msg.Content[0])
	writeServerData()
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Added album **"+msg.Content[0]+"** to album list.")
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
func delAlbumCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	checkChannelsMap(data)

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	if channel, ok := data.Channels[msg.ChannelID]; ok && len(channel.Albums) > 0 {
		data.Channels[msg.ChannelID].Albums = deleteAlbums(data.Channels[msg.ChannelID].Albums, msg.Content[0])
		writeServerData()
		result = fmt.Sprintf("Removed **%s** from the list of albums.", msg.Content[0])
	}

	_, _ = s.ChannelMessageSend(msg.ChannelID, result)

}

// Handles the image command.
// Checks whether an album is actually in cache before making an imgur API call.
func getImageCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	checkChannelsMap(data)

	var result string

	// Check there is server data for the given channel and check if there is at least 1 album.
	if channel, ok := data.Channels[msg.ChannelID]; !ok || len(channel.Albums) <= 0 {
		result = fmt.Sprintf("This channel does not have any albums, add an album using `%saddalbum <AlbumID> `.", data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
	} else {
		randomImageID := rand.Intn(time.Now().Nanosecond()) % len(channel.Albums)
		id := channel.Albums[randomImageID]
		data, ok := AlbumCache[id]
		if !ok {
			var err error
			data, _, err = imgClient.GetAlbumInfo(id)
			if err != nil {
				result = fmt.Sprintf("Something went wrong while trying to retrieve an image, maybe the Imgur API is down or **%s** is not an album?", id)
				log.Error(err)
				s.ChannelMessageSend(msg.ChannelID, result)
				return
			}
			AlbumCache[id] = data

		}
		rndImg := rand.Intn(time.Now().Nanosecond()) % len(data.Images)
		linkToImg := data.Images[rndImg].Link

		resultEmbed := NewEmbed().SetImage(linkToImg)

		_, _ = s.ChannelMessageSendEmbed(msg.ChannelID, resultEmbed.MessageEmbed)

	}
}

func getSauceCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	var url string

	if len(msg.Content) == 0 {
		amountOfMessages := 10
		res, err := findLastMessageWithAttachOrEmbed(s, msg, amountOfMessages)
		if err != nil {
			result = fmt.Sprintf("Unable to find an image to query.")
			_, _ = s.ChannelMessageSend(msg.ChannelID, result)
			return
		}
		url = res
	} else {
		url = msg.Content[0]
	}

	if !isValidUrl(url) {
		result = fmt.Sprintf("Could not parse command arguments to a URL")
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	// FIXME: Use meta-data rather than url to determine whether something is an image.
	hasExtension := false
	for _, ext := range extensions {
		if strings.HasSuffix(url, ext) {
			hasExtension = true
		}
	}

	if !hasExtension {
		result = "URL has a non-supported file extension."
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	sauceClient := saucenao.New(BotConfig.SauceNaoToken)

	sauceResult, err := sauceClient.FromURL(url)

	if err != nil || len(sauceResult.Data) == 0 {
		result = fmt.Sprintf("Something went wrong while contacting the Saucenao API." +
			"You could try sourcing your images manually at https://saucenao.com/")
		log.Errorf("%s \n %s", err, sauceResult)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	similarity, err := strconv.ParseFloat(sauceResult.Data[0].Header.Similarity, 32)
	if err != nil || similarity < 80.0 {
		result = fmt.Sprintf("No images found with a confidence over 80.")
	} else {
		result = fmt.Sprintf("Source found with %s%% confidence: <%s>", sauceResult.Data[0].Header.Similarity, sauceResult.Data[0].Data.ExtUrls[0])
	}
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)

}
