package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// This module consists of all the commands that have to do with the workings of GenBot.
// E.g. turning on and off certain modules and changing they key GenBot listens to.
// It is impossible to turn off this module.
type MetaModule struct {
	MetaCommands map[string]Command
	HelpString	string
}

func (cmd *MetaModule) setup() {
	cmd.MetaCommands = map[string]Command{}

	addCommanderCommand := Command{
		Name: "addcommander",
		Description: "This command adds a user as commander. Being a commander overwrites " +
			"the full permission system of the bot and will allow a user to execute any command.",
		Usage:      "Usage: `%saddcommander <@user>`",
		Permission: discordgo.PermissionManageServer,
		Execute:    addCommanderCommand,
	}
	cmd.MetaCommands["addcommander"] = addCommanderCommand

	delCommanderCommand := Command{
		Name:        "delcommander",
		Description: "This command removes a user as commander.",
		Usage:       "Usage: `%sdelcommander <@user>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     delCommanderCommand,
	}
	cmd.MetaCommands["delcommander"] = delCommanderCommand

	setKeyCommand := Command{
		Name:        "setkey",
		Description: "Changes the key the bot listens to",
		Usage:       "Usage: `%ssetkey <key>` (Note: key should be of length 1)",
		Permission:  discordgo.PermissionManageServer,
		Execute:     setKeyCommand,
	}
	cmd.MetaCommands["setkey"] = setKeyCommand

	helpCommand := Command{
		Name:        "help",
		Description: "Help provides you with more information about any default command.",
		Usage:       "Usage: `%shelp <command name>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     helpCommand,
	}
	cmd.MetaCommands["help"] = helpCommand

	commandListCommand := Command{
		Name:        "commands",
		Description: "Lists all default commands",
		Usage:       "Usage: `%scommandlist`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     commandListCommands,
	}
	cmd.MetaCommands["commandlist"] = commandListCommand
	cmd.MetaCommands["commands"] = commandListCommand
	cmd.MetaCommands["cmds"] = commandListCommand

	blockCommand := Command{
		Name:        "block",
		Description: "Blocks the given command for everyone",
		Usage:       "Usage: `%sblock <commandname>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     addBlockedCommand,
	}
	cmd.MetaCommands["block"] = blockCommand

	unblockCommand := Command{
		Name:        "unblock",
		Description: "Unblocks the given command",
		Usage:       "Usage: `%sunblock <commandname>`",
		Permission:  discordgo.PermissionManageServer,
		Execute:     delBlockedCommand,
	}
	cmd.MetaCommands["unblock"] = unblockCommand

}

func (cmd *MetaModule) retrieveCommands() map[string]Command {
	return cmd.MetaCommands
}

func (cmd *MetaModule) retrieveHelp() (moduleName string, info string) {
	moduleName = "Core Module"
	info = commandsToHelp(&cmd.HelpString, cmd.MetaCommands)
	return
}

func (cmd *MetaModule) execute(s *discordgo.Session, m *discordgo.MessageCreate, msg SentMessageData, serverData *ServerData) {
	if serverData.Key != msg.Key {
		return
	}

	if command, ok := cmd.MetaCommands[msg.Command]; ok {
		isCommander, ok := serverData.Commanders[m.Author.ID]
		perm, _ := s.UserChannelPermissions(msg.Author.ID, msg.ChannelID)

		// Check if user has the correct permission or whether user is a commander
		if perm&command.Permission == command.Permission || (ok && isCommander) {
			log.Infof("Executing command `%s` in server `%s` ", command.Name, serverData.ID)
			command.Execute(command, s, msg, serverData)
		} else {
			log.Infof("Use of %s command denied for permission level %d", command.Name, perm)
		}
	}
}

// Send a private message with the basic info of the bot
func helpCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	// Default help command
	if len(msg.Content) != 1 {
		commandUrl := HServer.getUrlFromID(data.ID)

		result := fmt.Sprintf("Heya, This is GenBot written by GenDoNL#8196. \n"+
			"For a list of built-in commands use **%scommandlist**. \n"+
			"For server specific commands, check out: %s \n\n"+
			"The source code of the bot can be found here: <https://github.com/GenDoNL/GenBot>", data.Key, commandUrl)

		s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	// Help for a certain command, loops over all modules to try and find the command the user wants help with.
	for _, module := range Modules {
		cmd := module.retrieveCommands()

		if command, ok := cmd[msg.Content[0]]; ok {
			tempUsage := fmt.Sprintf(command.Usage, msg.Key)
			result := fmt.Sprintf("***%s***\nDescription: %s\n\n%s\n", command.Name, command.Description, tempUsage)
			s.ChannelMessageSend(msg.ChannelID, result)
		}
	}
}

// Change the command key of the server.
// Key should be the first argument and only 1 char long.
func setKeyCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
	} else if len(msg.Content[0]) == 1 {
		data.Key = msg.Content[0]
		writeServerDataDB(data)
		result = fmt.Sprintf("Changed bot key to **%s**.", msg.Content[0])
	} else {
		result = fmt.Sprintf("Key should have a length of **1**")
	}
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

func commandListCommands(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	url := HServer.getUrlFromID(data.ID)
	result := fmt.Sprintf("use `%shelp <commandname>` for an in-depth description of commands.\n" +
		"The custom commands for this server can be found at %s\n", data.Key, url)

	embed := NewEmbed().
		SetDescription(result).
		SetAuthorFromUser(msg.Author).
		SetColorFromUser(s, msg.ChannelID, msg.Author)

	// Not using a for loop since we want the modules to be printed in a specific order.
	name, info := Modules["MetaModule"].retrieveHelp()
	embed.AddField(name, info)

	name, info = Modules["ModerationModule"].retrieveHelp()
	embed.AddField(name, info)

	name, info = Modules["CommandModule"].retrieveHelp()
	embed.AddField(name, info)

	name, info = Modules["ImageModule"].retrieveHelp()
	embed.AddField(name, info)

	s.ChannelMessageSendEmbed(msg.ChannelID, embed.MessageEmbed)
}

// Add a commander to the list of commanders.
// First argument should mention the user who should be added.
func addCommanderCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(data.Commanders) == 0 {
		data.Commanders = make(map[string]bool)
	}

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	userID, err := parseMention(msg.Content[0])

	if err != nil {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	data.Commanders[userID] = true
	writeServerDataDB(data)
	result = fmt.Sprintf("Added <@%s> as commander.", userID)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)

}

// Deletes a commander from the list of commanders.
// First argument should mention the user who should be deleted.
func delCommanderCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(data.Commanders) == 0 {
		data.Commanders = make(map[string]bool)
	}

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	userID, err := parseMention(msg.Content[0])

	if err != nil {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	data.Commanders[userID] = false
	writeServerDataDB(data)
	result = fmt.Sprintf("Removed <@%s> as commander.", userID)
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)
}

func addBlockedCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(data.BlockedCommands) == 0 {
		data.BlockedCommands = make(map[string]bool)
	}

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	if msg.Content[0] == "block" || msg.Content[0] == "unblock" {
		_, _ = s.ChannelMessageSend(msg.ChannelID, "Cannot block the block & unblock commands.")
		return
	}

	data.BlockedCommands[msg.Content[0]] = true
	writeServerDataDB(data)
	result = fmt.Sprintf("Blocked `%s`, use `%sunblock %s` to unblock the command.", msg.Content[0], data.Key, msg.Content[0])
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)

}

func delBlockedCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(data.BlockedCommands) == 0 {
		data.BlockedCommands = make(map[string]bool)
	}

	var result string

	if len(msg.Content) == 0 {
		result = fmt.Sprintf(command.Usage, data.Key)
		_, _ = s.ChannelMessageSend(msg.ChannelID, result)
		return
	}

	data.BlockedCommands[msg.Content[0]] = false
	writeServerDataDB(data)
	result = fmt.Sprintf("Unblocked `%s`.", msg.Content[0])
	_, _ = s.ChannelMessageSend(msg.ChannelID, result)

}
