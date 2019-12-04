package CoreModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"time"
)

func initPingCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "ping",
		description: "Returns pong.",
		usage:       "`%sping`",
		aliases:	 []string{},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*CoreModule).pingCommand,
	}
	return
}

func (c *CoreModule) pingCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {

	ping := s.HeartbeatLatency() / time.Millisecond

	response := Bot.NewEmbed().
		SetAuthorFromUser(m.Author).
		SetColorFromUser(s, m.ChannelID, m.Author).
		SetDescription(fmt.Sprintf("❤️ %dms", ping))

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, response.MessageEmbed)

	if err != nil {
		c.Bot.Log.Error(err)
		return
	}
}