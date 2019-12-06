package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
)


func initMeIrl() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "me_irl",
		description: "This commands sends your irl, given that you have an irl. You can add a me_irl by using `addme_irl`",
		usage:       "`%sme_irl`",
		aliases:	 []string{"meirl"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*CoreModule).meIrlCommand,
	}
	return
}

func (c *CoreModule) meIrlCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	var result string

	if res, ok := data.MeIrlData[m.Author.ID]; ok {
		result = res.Content
	} else {
		result = "You do not seem to have a me_irl."
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, result)
}