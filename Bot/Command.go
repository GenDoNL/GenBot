package Bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Command interface {
	Name() string
	Description() string
	Usage() string
	Permission() int
	Aliases() []string
}

func (b *Bot) CommandAllowed(c Command, cmdName string) bool {
	lowerName := strings.ToLower(cmdName)
	if lowerName == c.Name() {
		return true
	}

	for _, alias := range c.Aliases() {
		if lowerName == alias {
			return true
		}
	}
	return false
}

func (b *Bot) CanExecute(cmd Command, s *discordgo.Session, m *discordgo.MessageCreate, data *ServerData) bool  {
	perm, _ := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	isCommander, ok := data.Commanders[m.Author.ID]

	if perm&cmd.Permission() != cmd.Permission() && !(ok && isCommander) {
		Log.Infof("Use of %s command denied for permission level %d", cmd.Name(), perm)
		return false
	}
	return true
}

func (b *Bot) Usage(cmd Command, s *discordgo.Session, m *discordgo.MessageCreate, data *ServerData) *Embed {
	tempUsage := fmt.Sprintf(cmd.Usage(), data.Key)
	result := NewEmbed().
		SetAuthorFromUser(m.Author).
		SetColorFromUser(s, m.ChannelID, m.Author).
		SetTitle(fmt.Sprintf("%s%s", data.Key, cmd.Name())).
		SetDescription(cmd.Description()).
		AddField("Usage", tempUsage)
	return result
}

func stringInMap(str string, slice map[string]bool) (contains bool) {
	contains = false // Sanity check, default value should be false
	for k, v := range slice {
		if str == k {
			return v
		}
	}
	return
}

func (b *Bot) IsBlocked(cmd string, data *ServerData) (isBlocked bool) {
	return stringInMap(cmd, data.BlockedCommands)
}
