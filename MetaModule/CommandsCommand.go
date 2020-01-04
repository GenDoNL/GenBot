package MetaModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
)


func initCommandsCommand() (cc MetaCommand) {
	cc = MetaCommand{
		name:        "commands",
		description: "Show the default commands of GenBot",
		usage:       "`%scommands`",
		aliases:	 []string{"cmd", "command", "commandlist"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*MetaModule).commandsCommand,
	}
	return
}

func (c *MetaModule) commandsCommand(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	msg := fmt.Sprintf("use `%shelp <commandname>` for an in-depth description of commands.\n" +
		"The custom commands for this server can be found at %s\n", data.Key, c.Bot.Server.GetUrlFromID(data.ID))

	e := Bot.NewEmbed().
		SetDescription(msg).
		SetAuthorFromUser(m.Author).
		SetColorFromUser(s, m.ChannelID, m.Author)

	for _, module := range c.Bot.Modules {
		name, info := module.HelpFields()
		// Don't add empty fields, since that will cause a 400
		if name != "" && info != "" {
			e.AddField(name, info)
		}
	}

	ch, err := s.UserChannelCreate(m.Author.ID)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to send you a DM. A")
		Log.Error(err)
		return
	}

	_, err = s.ChannelMessageSendEmbed(ch.ID, e.MessageEmbed)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to send you a DM. B")
		Log.Error(err)
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

}