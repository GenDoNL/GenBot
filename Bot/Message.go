package Bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) FindLastMessageWithAttachOrEmbed(s *discordgo.Session, m *discordgo.MessageCreate, amount int) (result string, e error) {
	msgList, _ := s.ChannelMessages(m.ChannelID, amount, m.ID, "", "")

	copy(msgList[1:], msgList)
	// Prepend our actual message so we check this message for embeds as well
	msgList[0] = m.Message
	for _, x := range msgList {
		if len(x.Embeds) > 0 {
			if x.Embeds[0].URL != "" {
				result = x.Embeds[0].URL
				e = nil
				return
			}
		} else if len(x.Attachments) > 0 {
			result = x.Attachments[0].URL
			e = nil
			return
		}
	}

	result = ""
	e = errors.New("unable to find message with attachment or embed")
	return
}