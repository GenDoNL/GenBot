package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
	"time"
)

func initWhoIsCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "whois",
		description: "Sends member data of the mentioned user, or the message author if no-one is mentioned.",
		usage:       "`%swhois [user]`",
		aliases:	 []string{"whois", "who"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*CoreModule).whoIsCommand,
	}
	return
}

func (c *CoreModule) whoIsCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	targetName := ""
	if len(input) > 1 {
		targetName = input[1]
	}

	target := c.Bot.GetCommandTarget(s, m, data, targetName)

	memberData, err := s.GuildMember(data.ID, target.ID)
	if err != nil { return }

	response := Bot.NewEmbed().
		SetAuthorFromUser(target).
		SetColorFromUser(s, m.ChannelID, target).
		SetThumbnail(target.AvatarURL("256"))

	// Add nickname to message of the user has a nickname
	if memberData.Nick != "" {
		response.AddInlineField("Nickname", memberData.Nick, true)
	}

	// Set join and registration times.
	locale, _ := time.LoadLocation("UTC")
	joinTime, _ := memberData.JoinedAt.Parse()
	joinTimeParsed := joinTime.In(locale).Format(time.RFC1123)

	createTime := time.Unix(c.Bot.GetAccountCreationDate(target), 0).In(locale).Format(time.RFC1123)

	response.AddField("Registered", createTime).
		AddField("Joined", joinTimeParsed)

	// Add roles to the whois info
	roles := ""
	if len(memberData.Roles) != 0 {
		for _, roleId := range memberData.Roles {
			roles += "<@&" + roleId + "> "
		}

	} else {
		roles = "None"
	}
	response.AddField("Roles", roles)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, response.MessageEmbed)

	if err != nil {
		Log.Error(err)
		return
	}
}