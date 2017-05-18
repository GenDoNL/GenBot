package main

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

// Handles pruning of messages
func prune(s *discordgo.Session, msg MessageData, serverData *ServerData) {
	if len(msg.Content) == 0 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `"+serverData.Key+"prune <amount(1-100)> <user(optional)>`.")
		return
	}

	amount, err := strconv.Atoi(msg.Content[0])

	if err != nil || amount < 1 || amount > 100 {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Usage: `"+serverData.Key+"prune <amount(1-100)> <user(optional)>`.")
		return
	}

	// Retrieves 'amount' of messages before the command was issued.
	msgList, _ := s.ChannelMessages(msg.ChannelID, amount, msg.MessageID, "")

	var count = 0
	var msgID []string

	// Get the list of messages you want to remove.
	for _, x := range msgList {
		// Check if there was an user specified to be pruned, if so only prune that user.
		if len(msg.Mentions) == 0 || userInSlice(x.Author.ID, msg.Mentions) {
			count++
			msgID = append(msgID, x.ID)
		}
	}

	// Remove the messages
	s.ChannelMessagesBulkDelete(msg.ChannelID, msgID)

	// Send a conformation.
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Pruned **"+strconv.Itoa(count)+"** message(s).")

}
