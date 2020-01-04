package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)


func initAddMeIrl() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "addmeirl",
		description: "Add a me_irl, users can access this command themselves by using `me_irl`. " +
			"Other users can access the comment by using `<Nickname>_irl`",
		usage:       "`%saddme_irl <@User> <Nickname> <Content>`",
		aliases:	 []string{"addme_irl"},
		permission:  discordgo.PermissionManageServer,
		execute:     (*CoreModule).addMeIrlCommand,
	}
	return
}

func (c *CoreModule) addMeIrlCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 4)

	if len(input) < 4 || len(m.Mentions) == 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	targetId, err := c.Bot.ParseMention(input[1])
	if err != nil {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	nickname := input[2]
	cmdName := strings.ToLower(nickname) + "_irl"
	meIrl := Bot.MeIrlData{targetId, nickname, input[3]}

	err = data.CreateCustomCommand(cmdName, input[3])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while writing to the database, please try again later.")
		return
	}
	err = data.CreateMeIrl(meIrl)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while writing to the database, please try again later.")
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Successfully added me_irl command.")
	}
}