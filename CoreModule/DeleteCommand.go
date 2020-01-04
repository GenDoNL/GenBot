package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)

func initDeleteCommandCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "deletecommand",
		description: "Removes a custom command from this server.",
		usage:       "`%sdelcommand <Command_Name>`",
		aliases:	 []string{"delcommand", "rmcommand", "removecommand"},
		permission:  discordgo.PermissionManageServer,
		execute:     (*CoreModule).DeleteCommandCommand,
	}
	return
}

func (c *CoreModule) DeleteCommandCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 3)

	if len(input) < 2 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	deleted, err := data.DeleteCustomCommand(input[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while delete the command from the database, please try again later.")
		Log.Info(err)
		return
	}

	if !deleted {
		s.ChannelMessageSend(m.ChannelID, "Command has not been deleted, since no command with this name was found.")
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Successfully removed command.")
	}
}