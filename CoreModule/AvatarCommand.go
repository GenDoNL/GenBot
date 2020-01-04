package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)

func initAvatarCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "avatar",
		description: "Sends the full-size version of the mentioned user's avatar, or the message author if no-one is mentioned.",
		usage:       "`%savatar [user]`",
		aliases:	 []string{"avatar", "av"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*CoreModule).avatarCommand,
	}
	return
}

func (c *CoreModule) avatarCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	targetName := ""
	if len(input) > 1 {
		targetName = input[1]
	}

	target := c.Bot.GetCommandTarget(s, m, data, targetName)

	resultUrl := target.AvatarURL("256")
	response := Bot.NewEmbed().
		SetAuthorFromUser(target).
		SetColorFromUser(s, m.ChannelID, target).
		SetImage(resultUrl)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, response.MessageEmbed)

	if err != nil {
		Log.Error(err)
		return
	}
}