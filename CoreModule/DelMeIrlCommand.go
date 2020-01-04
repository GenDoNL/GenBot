package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
)


func initDelMeIrl() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "delmeirl",
		description: "Deletes a me_irl command",
		usage:       "`%sdelme_irl <@user>`",
		aliases:	 []string{"delme_irl", "rmme_irl", "rmmeirl"},
		permission:  discordgo.PermissionManageServer,
		execute:     (*CoreModule).delMeIrlCommand,
	}
	return
}

func (c *CoreModule) delMeIrlCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	if len(m.Mentions) == 0 {
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	target := m.Mentions[0].ID

	deleted, err := data.DeleteMeIrl(target)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Something went wrong while writing to the database, please try again later.")
		return
	}

	if !deleted {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No me_irl could be found for this user.")
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Successfully removed me_irl.")
	}
}