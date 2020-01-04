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
		SetColorFromUser(s, m.ChannelID, m.Author).
		SetTitle("🏓 Pong").
		SetDescription("Pinging...")

	// Benchmark round-trip time of sent message
	start := time.Now()
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, response.MessageEmbed)
	elapsed := time.Since(start) / time.Millisecond

	if err != nil {
		Log.Error(err)
		return
	}

	// Add the new data of the round-trip
	response.SetDescription("").
		AddInlineField("Bot", fmt.Sprintf("%dms", elapsed), true).
		AddInlineField("API", fmt.Sprintf("%dms", ping), true)

	_, err = s.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, response.MessageEmbed)
	if err != nil {
		Log.Error(err)
		return
	}
}