package MetaModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)


func initSetKeyCommand() (cc MetaCommand) {
	cc = MetaCommand{
		name:        "setkey",
		description: "Changes the key (or prefix) the bot listens to",
		usage:       "`%ssetkey <key (length 1)>`",
		aliases:	 []string{"setprefix", "setgenbotkey"},
		permission:  discordgo.PermissionManageServer,
		execute:     (*MetaModule).setKeyCommand,
	}
	return
}

func (c *MetaModule) setKeyCommand(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 3)

	if len(input) < 2 {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	if len(input[1]) != 1 {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	err := data.EditKey(input[1])

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while trying to update the key in the database, please try again later.")
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Successfully updated server key.")
	}
}