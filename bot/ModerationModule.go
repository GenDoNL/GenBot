package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"strconv"
)

type ModerationModule struct {
	ModerationCommands map[string]Command
}

func (cmd *ModerationModule) setup() {
	cmd.ModerationCommands = map[string]Command{}

	pruneCommand := Command{
		Name: "prune",
		Description: "This command prunes messages up to the provided amount. " +
			"If a user is mentioned, the command will only prune messages sent by this user.",
		Usage:      "Usage: `%sprune <amount(1-100)> <@user(optional, multiple)>`",
		Permission: discordgo.PermissionManageMessages,
		Execute:    pruneCommand,
	}
	cmd.ModerationCommands["prune"] = pruneCommand
}

func (cmd *ModerationModule) execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	serverData := getServerData(s, m.ChannelID)

	msg := parseMessage(m)

	if serverData.Key != msg.Key {
		return
	}

	if command, ok := cmd.ModerationCommands[msg.Command]; ok {
		isCommander, ok := serverData.Commanders[m.Author.ID]
		perm, _ := s.UserChannelPermissions(msg.Author.ID, msg.ChannelID)

		// Check if user has the correct permission or whether user is a commander
		if perm & command.Permission == command.Permission || (ok && isCommander) {
			log.Infof("Executing command: %s", command.Name)
			command.Execute(command, s, msg, serverData)
		} else {
			log.Infof("Use of %s command denied for permission level %d", command.Name, perm)
		}
	}
}


// Handles pruning of messages
func pruneCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	s.ChannelMessageDelete(msg.ChannelID, msg.MessageID)

	amount, err := strconv.Atoi(msg.Content[0])

	if err != nil || amount < 1 || amount > 100 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	// Retrieves a list of previously sent messages, up to `amount`
	msgList, _ := s.ChannelMessages(msg.ChannelID, amount, msg.MessageID, "", "")

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

	// Send a confirmation.
	result = fmt.Sprintf("Pruned **%s** message(s).", strconv.Itoa(count))
	_, _ = s.ChannelMessageSend(msg.ChannelID, "Pruned **"+strconv.Itoa(count)+"** message(s).")

}