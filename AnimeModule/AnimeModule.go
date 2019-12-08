package AnimeModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
)

type AnimeModule struct {
	Bot *Bot.Bot
	Commands []AnimeCommand
}

type AnimeCommand struct {
	name        string
	description string
	usage       string
	permission  int
	aliases     []string
	execute     func(*AnimeModule, AnimeCommand, *discordgo.Session, *discordgo.MessageCreate, *Bot.ServerData)
}

func New(bot *Bot.Bot) (c *AnimeModule) {
	c = &AnimeModule{Bot: bot}

	c.Commands = append(c.Commands, initMangaCommand())
	c.Commands = append(c.Commands, initAnimeCommand())
	c.Commands = append(c.Commands, initAniUserInfoCommand())

	return
}

func (c *AnimeModule) Execute(s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
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

func (c *AnimeModule) getCommand(cmdName string) (command AnimeCommand, found bool) {
	for _, cmd := range c.Commands {
		if c.Bot.CommandAllowed(cmd, cmdName) {
			command = cmd
			found = true
			return
		}
	}
	return
}

func (c *AnimeModule) CommandInfo(name string, data *Bot.ServerData) (response *Bot.Embed, found bool) {
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

func (c *AnimeModule) HelpFields() (title string, content string) {
	title = "Anime"
	for _, cmd := range c.Commands {
		if content == "" {
			content = cmd.Name()
		} else {
			content = fmt.Sprintf("%s, %s", content, cmd.Name())
		}
	}
	return
}

func (cc AnimeCommand) Name() string {
	return cc.name
}

func (cc AnimeCommand) Description() string {
	return cc.description
}

func (cc AnimeCommand) Usage() string {
	return cc.usage
}

func (cc AnimeCommand) Permission() int {
	return cc.permission
}

func (cc AnimeCommand) Aliases() []string {
	return cc.aliases
}



