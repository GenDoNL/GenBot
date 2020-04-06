package RSSModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"github.com/op/go-logging"
	"strings"
)

type RSSModule struct {
	Bot *Bot.Bot
	Commands []RSSCommand
}

type RSSCommand struct {
	name        string
	description string
	usage       string
	permission  int
	aliases     []string
	execute     func(*RSSModule, RSSCommand, *discordgo.Session, *discordgo.MessageCreate, *Bot.ServerData)
}
var (
	Log *logging.Logger
)

func New(bot *Bot.Bot, l *logging.Logger) (c *RSSModule) {
	c = &RSSModule{Bot: bot}
	Log = l

	go c.initPolling()
	c.Commands = append(c.Commands, initFollowRSSCommand())
	c.Commands = append(c.Commands, initUnfollowRSSCommand())

	return
}

func (c *RSSModule) Execute(s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
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

func (c *RSSModule) getCommand(cmdName string) (command RSSCommand, found bool) {
	for _, cmd := range c.Commands {
		if c.Bot.CommandAllowed(cmd, cmdName) {
			command = cmd
			found = true
			return
		}
	}
	return
}

func (c *RSSModule) CommandInfo(name string, data *Bot.ServerData) (response *Bot.Embed, found bool) {
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

func (c *RSSModule) HelpFields() (title string, content string) {
	title = "Core"
	for _, cmd := range c.Commands {
		if content == "" {
			content = cmd.Name()
		} else {
			content = fmt.Sprintf("%s, %s", content, cmd.Name())
		}
	}
	return
}

func (cc RSSCommand) Name() string {
	return cc.name
}

func (cc RSSCommand) Description() string {
	return cc.description
}

func (cc RSSCommand) Usage() string {
	return cc.usage
}

func (cc RSSCommand) Permission() int {
	return cc.permission
}

func (cc RSSCommand) Aliases() []string {
	return cc.aliases
}



