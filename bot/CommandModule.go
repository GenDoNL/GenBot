package main

import (
	"fmt"
	"github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type CommandModule struct {
	DefaultCommands map[string]Command
}

func (cmd *CommandModule) setup() {

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
		Permission:  discordgo.PermissionSendMessages,
		Execute:     rollCommand,
	}
	cmd.DefaultCommands["roll"] = rollCommand

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

	weatherCommand := Command{
		Name:        "weather",
		Description: "Send the current weather for the given location.",
		Usage:       "Usage: `%sweather <Location>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     weatherCommand,
	}
	cmd.DefaultCommands["weather"] = weatherCommand

	colorCommand := Command{
		Name:        "Color",
		Description: "Send the hex of the mentioned user, or the message author if no-one is mentioned.",
		Usage:       "Usage: `%scolor <@User(Optional)>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     colorCommand,
	}
	cmd.DefaultCommands["color"] = colorCommand
	cmd.DefaultCommands["colour"] = colorCommand
	cmd.DefaultCommands["clr"] = colorCommand

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
		writeServerDataDB(data)
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
	writeServerDataDB(data)
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

	writeServerDataDB(data)
	result = fmt.Sprintf("Removed the %sme_irl Command for <@%s>.", msg.Key, id)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

///// Default Commands /////

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

// Returns the avatar of a user
func avatarCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	target := getCommandTarget(s, msg, data)

	resultUrl := target.AvatarURL("256")
	result := NewEmbed().
		SetAuthorFromUser(target).
		SetColorFromUser(s, msg.ChannelID, target).
		SetImage(resultUrl)

	_, _ = s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)
}

func whoIsCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	target := getCommandTarget(s, msg, data)

	memberData, err := s.GuildMember(data.ID, target.ID)

	if err != nil {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Something went wrong while retrieving member data, please try again.")
		return
	}

	// Construct the base embed with user and avatar
	result := NewEmbed().
		SetAuthorFromUser(target).
		SetColorFromUser(s, msg.ChannelID, target).
		SetThumbnail(target.AvatarURL("256"))

	// Add nickname to message of the user has a nickname
	if memberData.Nick != "" {
		result.AddInlineField("Nickname", memberData.Nick, true)
	}

	// Set join and registration times.
	locale, _ := time.LoadLocation("UTC")
	joinTime, _ := memberData.JoinedAt.Parse()
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

func weatherCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	owm, err := openweathermap.NewCurrent("C", "EN", BotConfig.OwmToken)

	if err != nil {
		log.Error(err)
		return
	}

	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf(command.Usage, data.Key))
		return
	}

	err = owm.CurrentByName(strings.Join(msg.Content, " "))
	if err != nil {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Either the OpenWeatherMap API is down or you provided an invalid location.")
		return
	}

	fahr := owm.Main.Temp*9/5 + 32

	// Convert timezone data from Seconds to hours
	GmtOffset := owm.Timezone / 60 / 60
	localTime := time.Now().UTC().Add(time.Duration(GmtOffset) * time.Hour).Format("3:04PM, Monday") // Local time
	iconUrl := "http://openweathermap.org/img/wn/" + owm.Weather[0].Icon + "@2x.png"

	directionVal := int((owm.Wind.Deg / 22.5) + .5)
	directions := []string{"north", "north-northeast", "northeast", "east-northeast", "east", "east-southeast",
		"southeast", "south-southeast", "south", "south-southwest", "southwest", "west-southwest", "west", "west-northwest", "northwest", "north-northwest"}
	windDirection := directions[(directionVal % 16)]

	flag := fmt.Sprintf(":flag_%s:", strings.ToLower(owm.Sys.Country))

	result := NewEmbed().
		SetAuthorFromUser(msg.Author).
		SetColorFromUser(s, msg.ChannelID, msg.Author).
		SetThumbnail(iconUrl).
		SetTitle(fmt.Sprintf("Weather in **%s** %s at **%s**", owm.Name, flag, localTime)).
		AddField("Current Conditions:", fmt.Sprintf("**%s** at **%.1f°C** / **%.1f°F**",
			owm.Weather[0].Description, owm.Main.Temp, fahr)).
		AddInlineField("Humidity", fmt.Sprintf("%d%%", owm.Main.Humidity), true).
		AddInlineField("Wind", fmt.Sprintf("%.1f km/h from the %s ", owm.Wind.Speed*3.6, windDirection), true).
		SetFooter("Data provided by OpenWeatherMap", "http://f.gendo.moe/KlhvQJoD.png")

	_, _ = s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)

}

func colorCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	target := getCommandTarget(s, msg, data)

	color := s.State.UserColor(target.ID, msg.ChannelID)

	result := NewEmbed().
		SetAuthorFromUser(target).
		SetColor(color).
		SetDescription(fmt.Sprintf("\n#%x\n", color))

	_, _ = s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)
}
