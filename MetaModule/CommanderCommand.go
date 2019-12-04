package MetaModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
)


func initAddCommander() (cc MetaCommand) {
	cc = MetaCommand{
		name:        "addcommander",
		description: "This command adds a user as commander. Being a commander overwrites " +
				"the full permission system of the bot and will allow a user to execute any command.",
		usage:       "`%saddcommander <@user>`",
		aliases:	 []string{},
		permission:  discordgo.PermissionManageServer,
		execute:     (*MetaModule).addCommanderCommand,
	}
	return
}

func initDelCommander() (cc MetaCommand) {
	cc = MetaCommand{
		name:        "delcommander",
		description: "This command removes a user as commander.",
		usage:       "`%sdelcommander <@user>`",
		aliases:	 []string{"rmcommander"},
		permission:  discordgo.PermissionManageServer,
		execute:     (*MetaModule).delCommanderCommand,
	}
	return
}

func (c *MetaModule) addCommanderCommand(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	c.updateCommander(m, s, cmd, data, true)
}

func (c *MetaModule) delCommanderCommand(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	// This does leave the deleted commander in the database.
	// Allows one to easily update the value in the database if necessary.
	c.updateCommander(m, s, cmd, data, false)
}

func (c *MetaModule) updateCommander(m *discordgo.MessageCreate, s *discordgo.Session, cmd MetaCommand, data *Bot.ServerData, isCommander bool) {
	if len(m.Mentions) == 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	err := c.Bot.UpdateCommander(data.ID, m.Mentions[0].ID, isCommander)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while trying to update the key in the database, please try again later.")
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}