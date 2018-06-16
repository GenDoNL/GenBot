package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"strings"
)

type CommandModule struct {
	DefaultCommands map[string]Command
}

func (cmd *CommandModule) setup() {
	cmd.DefaultCommands = map[string]Command {}

	addCommand := Command {
		Name: "addcommand",
		Description: "Add a custom command",
		Usage: "Usage: `%saddcommand <command name> <response>`",
		Permission: discordgo.PermissionManageServer,
		Execute: addCommandCommand,
	}
	cmd.DefaultCommands["addcommand"] = addCommand
}

func (cmd *CommandModule) execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	serverData := getServerData(s, m.ChannelID)

	msg := parseMessage(m)

	if serverData.Key != msg.Key {
		return
	}

	if command, ok := cmd.DefaultCommands[msg.Command]; ok {
		isCommander, ok := serverData.Commanders[m.Author.ID]
		perm, _ := s.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
		if  perm & command.Permission == command.Permission || (ok && isCommander) {
			log.Infof("Executing command: %s", command.Name)
			command.Execute(command, s, msg, serverData)
		} else {
			log.Infof("Use of %s command denied for permission level %d", command.Name, perm)
		}
	}  else if cmd, ok := serverData.CustomCommands[msg.Command]; ok {
		s.ChannelMessageSend(m.ChannelID, cmd.Content)
	}
}

// Add a command to the list  of commands.
// First argument should be the name, the rest should be the content.
// Will overwrite existing commands!
func addCommandCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	checkCommandsMap(data)

	var result string

	if len(msg.Content) > 1 {
		createCommand(data, msg.Content[0], strings.Join(msg.Content[1:], " "))
		result = fmt.Sprintf("Added **%s** to the list of commands.", msg.Content[0])
	} else {
		result = fmt.Sprintf(command.Usage, data.Key)
	}

	s.ChannelMessageSend(msg.ChannelID, result)
}