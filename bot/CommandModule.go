package main

import (
	"github.com/GenDoNL/saucenao-go"
	"github.com/bwmarrin/discordgo"
	"github.com/koffeinsource/go-imgur"

	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CommandModule struct {
	DefaultCommands map[string]Command
}

func (cmd *CommandModule) setup() {
	imgClient = new(imgur.Client)
	imgClient.HTTPClient = new(http.Client)
	imgClient.Log = log
	imgClient.ImgurClientID = BotConfig.ImgurToken

	AlbumCache = make(map[string]*imgur.AlbumInfo)

	cmd.DefaultCommands = map[string]Command{}

	addCommand := Command{
		Name:        "addcommand",
		Description: "Adds a custom command",
		Usage:       "Usage: `%saddcommand <command name> <response>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     addCommandCommand,
	}
	cmd.DefaultCommands["addcommand"] = addCommand

	delCommand := Command{
		Name:        "addcommand",
		Description: "Removes a custom command",
		Usage:       "Usage: `%sdelcommand <command name>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     delCommandCommand,
	}
	cmd.DefaultCommands["delcommand"] = delCommand

	addMeIrlCommand := Command{
		Name:        "addme_irl",
		Description: "Add a me_irl",
		Usage:       "Usage: `%saddme_irl <@User> <Nickname> <Content>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     addMeIrlCommand,
	}
	cmd.DefaultCommands["addmeirl"] = addMeIrlCommand
	cmd.DefaultCommands["addme_irl"] = addMeIrlCommand

	delMeIrlCommand := Command{
		Name:        "delme_irl",
		Description: "Delete a me_irl",
		Usage:       "Usage: `%sdelme_irl <@User>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     delMeIrlCommand,
	}
	cmd.DefaultCommands["delmeirl"] = delMeIrlCommand
	cmd.DefaultCommands["delme_irl"] = delMeIrlCommand

	addAlbumCommand := Command{
		Name:        "addalbum",
		Description: "This command adds an album to this channel. Images can be retrieved at random using the `i` or `image` command.",
		Usage:       "Usage: `%saddalbum <imgur album id>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     addAlbumCommand,
	}
	cmd.DefaultCommands["addalbum"] = addAlbumCommand

	delAlbumCommand := Command{
		Name:        "delalbum",
		Description: "This command removes an album from this channel.",
		Usage:       "Usage: `%sdelalbum <imgur album id>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     delAlbumCommand,
	}
	cmd.DefaultCommands["delalbum"] = delAlbumCommand

	imageCommand := Command{
		Name:        "image",
		Description: "This command sends a random image from an album added by `addalbum`",
		Usage:       "Usage: `%si`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     getImageCommand,
	}
	cmd.DefaultCommands["i"] = imageCommand
	cmd.DefaultCommands["image"] = imageCommand

	meirlCommand := Command{
		Name:        "me_irl",
		Description: "This commands sends your irl, given that you have an irl. You can add an me_irl by using `addme_irl`",
		Usage:       "Usage: `%sme_irl`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     meIrlCommand,
	}
	cmd.DefaultCommands["meirl"] = meirlCommand
	cmd.DefaultCommands["me_irl"] = meirlCommand

	customCommandListCommand := Command{
		Name:        "customcommands",
		Description: "Lists all the server specific commands",
		Usage:       "Usage: `%scustomcommands`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     customCommandListCommand,
	}
	cmd.DefaultCommands["customcommands"] = customCommandListCommand
	cmd.DefaultCommands["servercommands"] = customCommandListCommand

	rollCommand := Command{
		Name:        "roll",
		Description: "Rolls a random number between 0 and the upper bound provided (0 and the upper bound included). Upper bound defaults to 100.",
		Usage:       "Usage: `%sroll <upper bound(optional)>`",
		Permission:  discordgo.PermissionManageMessages,
		Execute:     rollCommand,
	}
	cmd.DefaultCommands["roll"] = rollCommand

	sauceCommand := Command{
		Name:        "sauce",
		Description: "Provides the source of an image. A direct link to an image should be provided.",
		Usage:       "Usage: `%ssource <URL>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     getSauceCommand,
	}
	cmd.DefaultCommands["sauce"] = sauceCommand
	cmd.DefaultCommands["source"] = sauceCommand

	avatarCommand := Command{
		Name:        "avatar",
		Description: "Sends the full-size version of the mentioned user's avatar, or the message author if no-one is mentioned.",
		Usage:       "Usage: `%savatar <@User(Optional)>",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     avatarCommand,
	}

	cmd.DefaultCommands["avatar"] = avatarCommand
	cmd.DefaultCommands["av"] = avatarCommand

	whoIsCommand := Command{
		Name:        "whois",
		Description: "Sends member data of the mentioned user, or the message author if no-one is mentioned.",
		Usage:       "Usage: `%swhois <@User(Optional)>",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     whoIsCommand,
	}

	cmd.DefaultCommands["whois"] = whoIsCommand
	cmd.DefaultCommands["who"] = whoIsCommand
}

func (cmd *CommandModule) retrieveCommands() map[string]Command {
	return cmd.DefaultCommands
}

func (cmd *CommandModule) retrieveHelp() string {
	return ""
}

func (cmd *CommandModule) execute(s *discordgo.Session, m *discordgo.MessageCreate, msg SentMessageData, serverData *ServerData) {
	if serverData.Key != msg.Key {
		return
	}

	if command, ok := cmd.DefaultCommands[msg.Command]; ok {
		isCommander, ok := serverData.Commanders[m.Author.ID]
		perm, _ := s.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
		if perm&command.Permission == command.Permission || (ok && isCommander) {
			log.Infof("Executing command `%s` in server `%s` ", command.Name, serverData.ID)
			command.Execute(command, s, msg, serverData)
		} else {
			log.Infof("Use of %s command denied for permission level %d", command.Name, perm)
		}
	} else if cmd, ok := serverData.CustomCommands[msg.Command]; ok {
		log.Infof("Executing custom command `%s` in server `%s` ", cmd.Name, serverData.ID)

		s.ChannelMessageSend(m.ChannelID, cmd.Content)
	}
}

///// Admin Commands /////

// Add a command to the list  of commands.
// First argument should be the name, the rest should be the content.
// Will overwrite existing commands!
func addCommandCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	checkCommandsMap(data)

	var result string

	if len(msg.Content) > 1 {
		err := createCommand(data, msg.Content[0], strings.Join(msg.Content[1:], " "))

		if err != nil {
			result = fmt.Sprintf("Cannot add new commands that contain newlines in the name.")
		} else {
			result = fmt.Sprintf("Added **%s** to the list of commands.", msg.Content[0])
		}
	} else {
		result = fmt.Sprintf(command.Usage, data.Key)
	}

	s.ChannelMessageSend(msg.ChannelID, result)
}

// Removes a command from the list of commands
// First argument should be the name of the command that should be removed
// Nothing happens if command is unknown.
func delCommandCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	checkCommandsMap(data)

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
	} else if _, ok := data.CustomCommands[msg.Content[0]]; ok {
		delete(data.CustomCommands, msg.Content[0])
		writeServerData()
		result = fmt.Sprintf("Removed **%s** from the list of commands.", msg.Content[0])
	} else {
		result = fmt.Sprintf("The command **%s** has not been found.", msg.Content[0])
	}

	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

// Add a meIrlCommand command from a specific person
func addMeIrlCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if len(msg.Content) < 3 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	id, err := parseMention(msg.Content[0])

	if err != nil {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	nick := msg.Content[1]
	content := strings.Join(msg.Content[2:], " ")

	cmd := nick + "_irl"
	err = createCommand(data, cmd, content)
	if err != nil {
		result = fmt.Sprintf("Cannot add new commands that contain newlines in the name.")
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	if len(data.MeIrlData) == 0 {
		data.MeIrlData = make(map[string]*MeIrlData)
	}

	data.MeIrlData[id] = &MeIrlData{id, nick, content}
	writeServerData()
	result = fmt.Sprintf("Added %s%s", data.Key, cmd)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

// Remove a meIrlCommand command from a specific person.
func delMeIrlCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if len(msg.Content) < 1 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	id, err := parseMention(msg.Content[0])

	if err != nil {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	if meIrl, ok := data.MeIrlData[id]; ok {
		cmd := meIrl.Nickname + "_irl"
		delete(data.MeIrlData, id)
		if _, ok := data.CustomCommands[cmd]; ok {
			delete(data.CustomCommands, cmd)
		}
	}

	writeServerData()
	result = fmt.Sprintf("Removed the %sme_irl Command for <@%s>.", msg.Key, id)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
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

///// Default Commands /////

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

// Sends the me_irl command of an user or returns the default message.
func meIrlCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if res, ok := data.MeIrlData[msg.Author.ID]; ok {
		result = res.Content
	} else {
		result = "You do not seem to have a me_irl."
	}

	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

func rollCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var maxRoll int64 = 100

	if len(msg.Content) > 0 {
		newMaxRoll, err := strconv.ParseInt(msg.Content[0], 10, 64)
		if err == nil && newMaxRoll > 0 {
			maxRoll = newMaxRoll
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())

	result := strconv.FormatInt(rand.Int63n(maxRoll+1), 10)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

func customCommandListCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	url := HServer.updateServerCommands(data.ID, data)

	result := fmt.Sprintf("The full list of commands for this server can be found here: %s", url)

	s.ChannelMessageSend(msg.ChannelID, result)
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

// Returns the avatar of an user 
func avatarCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var target *discordgo.User

	if len(msg.Mentions) != 0 {
		target = msg.Mentions[0]
	} else {
		target = msg.Author
	}

	resultUrl := getAvatarFromUser(target)
	result := NewEmbed().SetAuthorFromUser(target).SetImage(resultUrl)

	_, _ = s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)
}

func whoIsCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var target *discordgo.User

	if len(msg.Mentions) != 0 {
		target = msg.Mentions[0]
	} else {
		target = msg.Author
	}

	memberData, err := s.GuildMember(data.ID, target.ID)

	if err != nil {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Something went wrong while retrieving member data, please try again.")
		return
	}


	// Construct the base embed with user and avatar
	result := NewEmbed().
		SetAuthorFromUser(target).
		SetThumbnail(getAvatarFromUser(target))

	// Add nickname to message of the user has a nickname
	if memberData.Nick != "" {
		result.AddField("Nickname", memberData.Nick)
	}

	// Set join and registration times.
	locale, _ := time.LoadLocation("UTC")
	joinTime, _ := time.Parse(time.RFC3339, memberData.JoinedAt)
	longTime := joinTime.In(locale).Format(time.RFC1123)

	createTime := time.Unix(getAccountCreationDate(target), 0).In(locale).Format(time.RFC1123)
	result.AddField("Registered", createTime).AddField("Joined", longTime)

	// Add roles to the whois info
	roles := ""
	if len(memberData.Roles) != 0 {
		for _, roleId := range memberData.Roles {
			roles += "<@&" + roleId + "> "
		}

	} else {
		roles = "None"
	}
	result.AddField("Roles", roles)
	_, _ = s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)

}
