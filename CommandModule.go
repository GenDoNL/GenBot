package main

import (
	"github.com/GenDoNL/saucenao-go"
	"github.com/bwmarrin/discordgo"
	"github.com/koffeinsource/go-imgur"

	"fmt"
	"math/rand"
	"net/http"
	"sort"
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

	setKeyCommand := Command{
		Name:        "setkey",
		Description: "Changes the key the bot listens to",
		Usage:       "Usage: `%ssetkey <key>` (Note: key should be of length 1)",
		Permission:  discordgo.PermissionManageServer,
		Execute:     setKeyCommand,
	}
	cmd.DefaultCommands["setkey"] = setKeyCommand

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

	lockCommand := Command{
		Name:        "lock",
		Description: "This command disallows the `everyone` role to speak in the current channel.",
		Usage:       "Usage: `%slock`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     lockChannelCommand,
	}
	cmd.DefaultCommands["lock"] = lockCommand

	unlockCommand := Command{
		Name:        "unlock",
		Description: "This command allows the `everyone` role to speak in the current channel.",
		Usage:       "Usage: `%sunlock`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     unlockChannelCommand,
	}
	cmd.DefaultCommands["unlock"] = unlockCommand

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

	addCommanderCommand := Command{
		Name: "addcommander",
		Description: "This command adds a user as commander. Being a commander overwrites " +
			"the full permission system of the boot and will allow a user to execute any command.",
		Usage:      "Usage: `%saddcommander <@user>`",
		Permission: discordgo.PermissionManageServer,
		Execute:    addCommanderCommand,
	}
	cmd.DefaultCommands["addcommander"] = addCommanderCommand

	delCommanderCommand := Command{
		Name:        "delcommander",
		Description: "This command removes a user as commander.",
		Usage:       "Usage: `%sdelcommander <@user>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     delCommanderCommand,
	}
	cmd.DefaultCommands["delcommander"] = delCommanderCommand

	pruneCommand := Command{
		Name: "prune",
		Description: "This command prunes messages up to the provided amount. " +
			"If a user is mentioned, the command will only prune messages sent by this user.",
		Usage:      "Usage: `%sprune <amount(1-100)> <@user(optional, multiple)>`",
		Permission: discordgo.PermissionManageMessages,
		Execute:    pruneCommand,
	}
	cmd.DefaultCommands["prune"] = pruneCommand

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

	commandListCommand := Command{
		Name:        "commandlist",
		Description: "Lists all default commands",
		Usage:       "Usage: `%scommandlist`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     commandListCommands,
	}
	cmd.DefaultCommands["commandlist"] = commandListCommand
	cmd.DefaultCommands["commands"] = commandListCommand

	helpCommand := Command{
		Name:        "help",
		Description: "Help provides you with more information about any default command.",
		Usage:       "Usage: `%shelp <command name>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     helpCommand,
	}
	cmd.DefaultCommands["help"] = helpCommand

	rolLCommand := Command{
		Name:        "roll",
		Description: "Rolls a random number between 0 and the upper bound provided (0 and the upper bound included). Upper bound defaults to 100.",
		Usage:       "Usage: `%sroll <upper bound(optional)>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     rollCommand,
	}
	cmd.DefaultCommands["roll"] = rolLCommand

	sauceCommand := Command{
		Name:        "sauce",
		Description: "Provides the source of an image. A direct link to an image should be provided.",
		Usage:       "Usage: `%ssource <URL>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     getSauceCommand,
	}
	cmd.DefaultCommands["sauce"] = sauceCommand
	cmd.DefaultCommands["source"] = sauceCommand

}

func (cmd *CommandModule) execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	serverData := getServerData(s, m.ChannelID)

	msg := parseMessage(m)

	if serverData.Key != msg.Key {
		return
	}

	if command, ok := cmd.DefaultCommands[msg.Command]; ok {
		isCommander, ok := serverData.Commanders[m.Author.ID]
		perm, _ := s.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
		if perm&command.Permission == command.Permission || (ok && isCommander) {
			log.Infof("Executing command: %s", command.Name)
			command.Execute(command, s, msg, serverData)
		} else {
			log.Infof("Use of %s command denied for permission level %d", command.Name, perm)
		}
	} else if cmd, ok := serverData.CustomCommands[msg.Command]; ok {
		log.Infof("Executing command: %s, from server: %s", cmd.Name, serverData.ID)

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
		createCommand(data, msg.Content[0], strings.Join(msg.Content[1:], " "))
		result = fmt.Sprintf("Added **%s** to the list of commands.", msg.Content[0])
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

// Change the command key of the server.
// Key should be the first argument and only 1 char long.
func setKeyCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
	} else if len(msg.Content[0]) == 1 {
		data.Key = msg.Content[0]
		writeServerData()
		result = fmt.Sprintf("Changed bot key to **%s**.", msg.Content[0])
	} else {
		result = fmt.Sprintf("Key should have a length of **1**")
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
	createCommand(data, cmd, content)

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

// Lock a channel, so the @everyone role won't be able to talk in the channel.
func lockChannelCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {

	//get channel object
	ch, _ := s.Channel(msg.ChannelID)
	sv, _ := s.Guild(data.ID)

	everyonePerms, err := getRolePermissionsByName(ch, sv, "@everyone")

	//deny sending messages and update it
	err = s.ChannelPermissionSet(ch.ID, everyonePerms.ID, everyonePerms.Type, everyonePerms.Allow&^0x800, everyonePerms.Deny|0x800)

	var result string
	if err != nil {
		result = fmt.Sprintf("Unable to unlock channel, do I have the permissions to manage roles?")
		log.Errorf("Error unlocking channel: %s", err)
	} else {
		result = fmt.Sprintf("This channel has been locked.")
		log.Infof("Locked channel: %s, in server: %s", msg.ChannelID, data.ID)
	}
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

// Unlock the channel, so the @everyone role will be allowed to talk again.
func unlockChannelCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	//get channel object
	ch, _ := s.Channel(msg.ChannelID)
	sv, _ := s.Guild(data.ID)

	everyonePerms, err := getRolePermissionsByName(ch, sv, "@everyone")

	//deny sending messages and update it
	err = s.ChannelPermissionSet(ch.ID, everyonePerms.ID, everyonePerms.Type, everyonePerms.Allow|0x800, everyonePerms.Deny&^0x800)

	var result string

	if err != nil {
		result = fmt.Sprintf("Unable to unlock channel, do I have the permissions to manage roles?")
		log.Errorf("Error unlocking channel: %s", err)
	} else {
		result = fmt.Sprintf("This channel has been unlocked.")
		log.Infof("Unlocked channel: %s, in server: %s", msg.ChannelID, data.ID)
	}

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

// Deletes a commander from the list of commanders.
// First argument should mention the user who should be deleted.
func delCommanderCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(data.Commanders) == 0 {
		data.Commanders = make(map[string]bool)
	}

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	userID, err := parseMention(msg.Content[0])

	if err != nil {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	data.Commanders[userID] = false
	writeServerData()
	result = fmt.Sprintf("Removed <@%s> as commander.", userID)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

// Add a commander to the list of commanders.
// First argument should mention the user who should be added.
func addCommanderCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(data.Commanders) == 0 {
		data.Commanders = make(map[string]bool)
	}

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	userID, err := parseMention(msg.Content[0])

	if err != nil {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	data.Commanders[userID] = true
	writeServerData()
	result = fmt.Sprintf("Added <@%s> as commander.", userID)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)

}

///// Moderator Commands /////

// Handles pruning of messages
func pruneCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	amount, err := strconv.Atoi(msg.Content[0])

	if err != nil || amount < 1 || amount > 100 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	// Retrieves a list of previously sent messages, up to `amount`
	msgList, _ := s.ChannelMessages(msg.ChannelID, amount, msg.MessageID, "", "")

	var count = 0
	var msgID []string

	// Get the list of messages you want to remove.
	for _, x := range msgList {
		// Check if there was an user specified to be pruned, if so only prune that user.
		if len(msg.Mentions) == 0 || userInSlice(x.Author.ID, msg.Mentions) {
			count++
			msgID = append(msgID, x.ID)
		}
	}

	// Remove the messages
	s.ChannelMessagesBulkDelete(msg.ChannelID, msgID)

	// Send a conformation.
	result = fmt.Sprintf("Pruned **%s** message(s).", strconv.Itoa(count))
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Pruned **"+strconv.Itoa(count)+"** message(s).")

}

///// Default Commands /////

// Handles the i or image command.
// Checks whether an album is actually in cache before making an imgur API call.
func getImageCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	checkChannelsMap(data)

	var result string

	// Check there is server data for the given channel and check if there is at least 1 album.
	if channel, ok := data.Channels[msg.ChannelID]; ok && len(channel.Albums) > 0 {
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
				result = fmt.Sprintf("Something went wrong while trying to retrieve an image, maybe the Imgur API is down or **%s** is not an album?", id)
				log.Error(err)
				s.ChannelMessageSend(msg.ChannelID, result)
				return
			}
			AlbumCache[id] = data

		}
		// Get a random image from the album and get the link of said image.
		rndImg := rand.Intn(time.Now().Nanosecond()) % len(data.Images)
		linkToImg := data.Images[rndImg].Link

		s.ChannelMessageSend(msg.ChannelID, linkToImg)
	} else {
		result = fmt.Sprintf("This channel does not have any albums, add an album using `%saddalbum <AlbumID> `.", data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
	}
}

// Retrieves the meIrlCommand of a given user.
func meIrlCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if res, ok := data.MeIrlData[msg.Author.ID]; ok {
		result = res.Content
	} else {
		result = "I sincerely apologize, you do not seem to have a me_irl."
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

// Send a private message with the basic info of the bot
func helpCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(msg.Content) == 0 {
		result := fmt.Sprintf("Heya, This is GenBot written by GenDoNL. \n"+
			"For the list of commands use **%scommandlist**. \n\n"+
			"The source code of the bot can be found here: https://github.com/GenDoNL/GenBot", data.Key)

		s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	if command, ok := CmdModule.DefaultCommands[msg.Content[0]]; ok {
		tempUsage := fmt.Sprintf(command.Usage, msg.Key)
		result := fmt.Sprintf("***%s***\nDescription: %s\n\n%s\n", command.Name, command.Description, tempUsage)
		s.ChannelMessageSend(msg.ChannelID, result)
	}
}

// Get all the message
func commandListCommands(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	result := fmt.Sprintf("This is a list of all default commands, use `%shelp <commandname>` for an in-depth description.\n\n", data.Key)

	var keys []string

	for k := range CmdModule.DefaultCommands {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, v := range keys {
		if !strings.Contains(result, fmt.Sprintf(" %s,", v)) {
			result = fmt.Sprintf("%s%s, ", result, v)

		}
	}

	if len(result) > 1950 {
		result = result[0:1950] + "...truncated"
	}

	s.ChannelMessageSend(msg.ChannelID, result)
}

func getSauceCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	if !isValidUrl(msg.Content[0]) {
		result = fmt.Sprintf("Could not parse command arguments to a URL")
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	hasExtension := false
	for _, ext := range extensions {
		if strings.HasSuffix(msg.Content[0], ext) {
			hasExtension = true
		}
	}

	if !hasExtension {
		result = "URL has a non-supported file extension."
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	sauceClient := saucenao.New(BotConfig.SauceNaoToken)

	sauceResult, err := sauceClient.FromURL(msg.Content[0])

	if err != nil || len(sauceResult.Data) == 0 {
		result = fmt.Sprintf("Something went wrong while contacting the Saucenao API." +
			"You could try sourcing your images manually at https://saucenao.com/")
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	similarity, err := strconv.ParseFloat(sauceResult.Data[0].Header.Similarity, 32)
	if err != nil || similarity < 80.0 {
		result = fmt.Sprintf("No images found with similarity over 80.")
	} else {
		result = fmt.Sprintf("Source found with %v%% similarity: <%s>", similarity, sauceResult.Data[0].Data.ExtUrls[0])
	}
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)

}
