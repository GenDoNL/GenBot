package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)

func initAddCommandCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "addcommand",
		description: "Adds a custom command to this server.",
		usage:       "`%saddcommand <Command_Name> <Content>`",
		aliases:	 []string{},
		permission:  discordgo.PermissionManageServer,
		execute:     (*CoreModule).AddCommandCommand,
	}
	return
}

func (c *CoreModule) AddCommandCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 3)

	if len(input) < 3 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	err := data.CreateCustomCommand(input[1], input[2])

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while writing the command to the database, please try again later.")
		Log.Info(err)
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Successfully added command.")
	}
}