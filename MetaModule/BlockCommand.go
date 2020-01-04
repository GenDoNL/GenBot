package MetaModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)


func initBlockCommand() (cc MetaCommand) {
	cc = MetaCommand{
		name:        "block",
		description: "Blocks a command, note that it does not block aliases. Meant to block conflicting command names.",
		usage:       "`%sblock <command_name>`",
		aliases:	 []string{},
		permission:  discordgo.PermissionManageServer,
		execute:     (*MetaModule).blockCommand,
	}
	return
}

func initUnblockCommand() (cc MetaCommand) {
	cc = MetaCommand{
		name:        "unblock",
		description: "Unblocks the given command.",
		usage:       "`%sblock <command_name>`",
		aliases:	 []string{},
		permission:  discordgo.PermissionManageServer,
		execute:     (*MetaModule).unblockCommand,
	}
	return
}

func (c *MetaModule) blockCommand(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	c.updateCommand(m, s, cmd, data, true)
}

func (c *MetaModule) unblockCommand(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	c.updateCommand(m, s, cmd, data, false)
}

func (c *MetaModule) updateCommand(m *discordgo.MessageCreate, s *discordgo.Session, cmd MetaCommand, data *Bot.ServerData, isBlocked bool) {
	input := strings.SplitN(m.Content, " ", 3)
	if len(input) < 2 {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	if input[1] == "block" || input[1] == "unblock" {
		s.ChannelMessageSend(m.ChannelID, "Cannot block the block and unblock commands.")
		return
	}

	err := data.UpdateBlockedCommand(input[1], isBlocked)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while trying to update the key in the database, please try again later.")
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Successfully updated blocked command.")
	}}