package MetaModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)


func initHelpCommand() (cc MetaCommand) {
	cc = MetaCommand{
		name:        "help",
		description: "Shows the help information for GenBot if no argument is given. " +
			"If an argument is given it will provide usage info about that command",
		usage:       "`%shelp [command_name]`",
		aliases:	 []string{},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*MetaModule).helpCommand,
	}
	return
}

func (c *MetaModule) helpCommand(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 3)

	if len(input) > 1 { // Case: We want info on specific command
		c.commandInfo(cmd, s, m, data, input[1])
	} else { // Case: We want the default help message
		c.basicInfo(cmd, s, m, data)
	}
}

func (c *MetaModule) basicInfo(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	result := fmt.Sprintf("This is GenBot written by Gen#8196. \n"+
				"For a list of built-in commands use **%scommands**. \n"+
				"For server specific commands, check out: %s \n\n"+
				"The source code of the bot can be found here: <https://github.com/GenDoNL/GenBot>", data.Key, c.Bot.Server.GetUrlFromID(data.ID))

	me, err := s.User("@me")

	if err != nil {
		return
	}

	e := Bot.NewEmbed().
		SetAuthorFromUser(m.Author).
		SetColorFromUser(s, m.ChannelID, m.Author).
		SetThumbnail(me.AvatarURL("256")).
		SetTitle("GenBot - Help").
		SetURL("https://github.com/GenDoNL/GenBot").
		SetDescription(result)

	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)

}

func (c *MetaModule) commandInfo(cmd MetaCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData, cmdName string) {
	for _, module := range c.Bot.Modules {
		embed, found := module.CommandInfo(cmdName, data)
		if found {
			embed.SetAuthorFromUser(m.Author).SetColorFromUser(s, m.ChannelID, m.Author)
			s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
			return
		}
	}
	s.ChannelMessageSend(m.ChannelID, "No command with this name has been found ")
}