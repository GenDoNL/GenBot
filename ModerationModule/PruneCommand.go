package ModerationModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strconv"
	"strings"
)

func initPruneCommand() (cc ModerationCommand) {
	cc = ModerationCommand{
		name:        "prune",
		description: "This command prunes messages up to the provided amount. " +
			"If any users are mentioned, the command will only prune messages sent by these users.",
		usage:       "`%sprune <amount(1-100)> [@user(multiple)]`",
		aliases:	 []string{"purge"},
		permission:  discordgo.PermissionManageMessages,
		execute:     (*ModerationModule).pruneCommand,
	}
	return
}

func (c *ModerationModule) pruneCommand(cmd ModerationCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 3)
	if len(input) < 2 {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	s.ChannelMessageDelete(m.ChannelID, m.ID)

	amount, err := strconv.Atoi(input[1])
	if err != nil || amount < 1 || amount > 100 {
		s.ChannelMessageSendEmbed(m.ChannelID, c.Bot.Usage(cmd, s, m, data).MessageEmbed)
		return
	}

	msgList, _ := s.ChannelMessages(m.ChannelID, amount, m.ID, "", "")

	var count = 0
	var msgID []string

	// Get the list of messages you want to remove.
	for _, x := range msgList {
		// Check if there was an user specified to be pruned, if so only prune that user.
		if len(m.Mentions) == 0 || userInSlice(x.Author.ID, m.Mentions) {
			count++
			msgID = append(msgID, x.ID)
		}
	}
	s.ChannelMessagesBulkDelete(m.ChannelID, msgID)

	result := fmt.Sprintf("Pruned **%s** message(s).", strconv.Itoa(count))
	_, _ = s.ChannelMessageSend(m.ChannelID, result)
}

func userInSlice(a string, list []*discordgo.User) bool {
	for _, b := range list {
		if b.ID == a {
			return true
		}
	}
	return false
}