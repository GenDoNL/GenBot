package MetaModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)

type MetaModule struct {
	Bot *Bot.Bot
	Commands []MetaCommand
}

type MetaCommand struct {
	name        string
	description string
	usage       string
	permission  int
	aliases     []string
	execute     func(*MetaModule, MetaCommand, *discordgo.Session, *discordgo.MessageCreate, *Bot.ServerData)
}

func New(bot *Bot.Bot) (c *MetaModule) {
	c = &MetaModule{Bot: bot}

	c.Commands = append(c.Commands, initHelpCommand())
	c.Commands = append(c.Commands, initCommandsCommand())
	c.Commands = append(c.Commands, initSetKeyCommand())
	c.Commands = append(c.Commands, initAddCommander())
	c.Commands = append(c.Commands, initDelCommander())
	c.Commands = append(c.Commands, initBlockCommand())
	c.Commands = append(c.Commands, initUnblockCommand())

	return
}

func (c *MetaModule) Execute(s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	cmdName := strings.ToLower(strings.SplitN(m.Content, " ", 2)[0])
	if cmdName[:1] != data.Key {
		return
	}

	command, found := c.getCommand(cmdName[1:])
	if !found {
		return
	}

	if !c.Bot.CanExecute(command, s, m, data) || c.Bot.IsBlocked(cmdName[1:], data) {
		return
	}

	c.Bot.Log.Infof("Executing command `%s` in server `%s` ", command.Name(), data.ID)
	command.execute(c, command, s, m, data)
}

func (c *MetaModule) getCommand(cmdName string) (command MetaCommand, found bool) {
	for _, cmd := range c.Commands {
		if c.Bot.CommandAllowed(cmd, cmdName) {
			command = cmd
			found = true
			return
		}
	}
	return
}

func (c *MetaModule) CommandInfo(name string, data *Bot.ServerData) (response *Bot.Embed, found bool) {
	cmd, found := c.getCommand(name)
	if !found { return }

	tempUsage := fmt.Sprintf(cmd.Usage(), data.Key)
	cmdName := fmt.Sprintf("%s%s", data.Key, cmd.Name())
	response = Bot.NewEmbed().
		SetTitle(fmt.Sprintf(cmdName)).
		SetDescription(cmd.Description()).
		AddField("Usage", tempUsage)
	return
}

func (c *MetaModule) HelpFields() (title string, content string) {
	title = "Meta"
	for _, cmd := range c.Commands {
		if content == "" {
			content = cmd.Name()
		} else {
			content = fmt.Sprintf("%s, %s", content, cmd.Name())
		}
	}
	return
}

func (cc MetaCommand) Name() string {
	return cc.name
}

func (cc MetaCommand) Description() string {
	return cc.description
}

func (cc MetaCommand) Usage() string {
	return cc.usage
}

func (cc MetaCommand) Permission() int {
	return cc.permission
}

func (cc MetaCommand) Aliases() []string {
	return cc.aliases
}



