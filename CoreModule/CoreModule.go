package CoreModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"github.com/op/go-logging"
	"strings"
)

type CoreModule struct {
	Bot *Bot.Bot
	Commands []CoreCommand
}

type CoreCommand struct {
	name        string
	description string
	usage       string
	permission  int
	aliases     []string
	execute     func(*CoreModule, CoreCommand, *discordgo.Session, *discordgo.MessageCreate, *Bot.ServerData)
}
var (
	Log *logging.Logger
)

func New(bot *Bot.Bot, l *logging.Logger) (c *CoreModule) {
	c = &CoreModule{Bot: bot}
	Log = l

	c.Commands = append(c.Commands, initPingCommand())
	c.Commands = append(c.Commands, initAddCommandCommand())
	c.Commands = append(c.Commands, initDeleteCommandCommand())
	c.Commands = append(c.Commands, initAvatarCommand())
	c.Commands = append(c.Commands, initWhoIsCommand())
	c.Commands = append(c.Commands, initColorCommand())
	c.Commands = append(c.Commands, initRollCommand())
	c.Commands = append(c.Commands, initWeatherCommand())
	c.Commands = append(c.Commands, initAddMeIrl())
	c.Commands = append(c.Commands, initDelMeIrl())
	c.Commands = append(c.Commands, initMeIrl())
	c.Commands = append(c.Commands, initSourceCommand())

	return
}

func (c *CoreModule) Execute(s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	cmdName := strings.ToLower(strings.SplitN(m.Content, " ", 2)[0])
	if cmdName[:1] != data.Key {
		return
	}

	command, found := c.getCommand(cmdName[1:])
	if !found {
		c.executeCustom(s, m, data, cmdName[1:])
		return
	}

	if !c.Bot.CanExecute(command, s, m, data) || c.Bot.IsBlocked(cmdName[1:], data) {
		return
	}

	Log.Infof("Executing command `%s` in server `%s` ", command.Name(), data.ID)
	command.execute(c, command, s, m, data)
}

func (c *CoreModule) executeCustom(s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData, cmdName string) {
	if cmd, ok := data.CustomCommands[cmdName]; ok {
		Log.Infof("Executing custom command `%s` in server `%s` ", cmd.Name, data.ID)

		s.ChannelMessageSend(m.ChannelID, cmd.Content)
	}
}

func (c *CoreModule) getCommand(cmdName string) (command CoreCommand, found bool) {
	for _, cmd := range c.Commands {
		if c.Bot.CommandAllowed(cmd, cmdName) {
			command = cmd
			found = true
			return
		}
	}
	return
}

func (c *CoreModule) CommandInfo(name string, data *Bot.ServerData) (response *Bot.Embed, found bool) {
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

func (c *CoreModule) HelpFields() (title string, content string) {
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

func (cc CoreCommand) Name() string {
	return cc.name
}

func (cc CoreCommand) Description() string {
	return cc.description
}

func (cc CoreCommand) Usage() string {
	return cc.usage
}

func (cc CoreCommand) Permission() int {
	return cc.permission
}

func (cc CoreCommand) Aliases() []string {
	return cc.aliases
}



