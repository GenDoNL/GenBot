package ModerationModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"github.com/op/go-logging"
	"strings"
)

type ModerationModule struct {
	Bot *Bot.Bot
	Commands []ModerationCommand
}

type ModerationCommand struct {
	name        string
	description string
	usage       string
	permission  int
	aliases     []string
	execute     func(*ModerationModule, ModerationCommand, *discordgo.Session, *discordgo.MessageCreate, *Bot.ServerData)
}

var (
	Log *logging.Logger
)

func New(bot *Bot.Bot, l *logging.Logger) (c *ModerationModule) {
	c = &ModerationModule{Bot: bot}
	Log = l

	c.Commands = append(c.Commands, initPruneCommand())
	c.Commands = append(c.Commands, initLockCommand())
	c.Commands = append(c.Commands, initUnlockCommand())

	return
}

func (c *ModerationModule) Execute(s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
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

	Log.Infof("Executing command `%s` in server `%s` ", command.Name(), data.ID)
	command.execute(c, command, s, m, data)
}

func (c *ModerationModule) getCommand(cmdName string) (command ModerationCommand, found bool) {
	for _, cmd := range c.Commands {
		if c.Bot.CommandAllowed(cmd, cmdName) {
			command = cmd
			found = true
			return
		}
	}
	return
}

func (c *ModerationModule) CommandInfo(name string, data *Bot.ServerData) (response *Bot.Embed, found bool) {
	cmd, found := c.getCommand(name)
	if !found { return }

	tempUsage := fmt.Sprintf(cmd.Usage(), data.Key)
	cmdName := fmt.Sprintf("%s%s", data.Key, cmd.Name())
	response =  Bot.NewEmbed().
		SetTitle(fmt.Sprintf(cmdName)).
		SetDescription(cmd.Description()).
		AddField("Usage", tempUsage)
	return
}

func (c *ModerationModule) HelpFields() (title string, content string) {
	title = "Moderation"
	for _, cmd := range c.Commands {
		if content == "" {
			content = cmd.Name()
		} else {
			content = fmt.Sprintf("%s, %s", content, cmd.Name())
		}
	}
	return
}

func (cc ModerationCommand) Name() string {
	return cc.name
}

func (cc ModerationCommand) Description() string {
	return cc.description
}

func (cc ModerationCommand) Usage() string {
	return cc.usage
}

func (cc ModerationCommand) Permission() int {
	return cc.permission
}

func (cc ModerationCommand) Aliases() []string {
	return cc.aliases
}



