package AnimeModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	anilistgo "github.com/gendonl/anilist-go"
	"github.com/gendonl/genbot/Bot"
	"github.com/op/go-logging"
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

var (
	Log *logging.Logger
)

func New(bot *Bot.Bot, l *logging.Logger) (c *AnimeModule) {
	c = &AnimeModule{Bot: bot}
	Log = l

	c.Commands = append(c.Commands, initMangaCommand())
	c.Commands = append(c.Commands, initAnimeCommand())
	c.Commands = append(c.Commands, initAniUserInfoCommand())
	c.Commands = append(c.Commands, initAniRecentCommand())
	c.Commands = append(c.Commands, initSetUserCommand())

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

	Log.Infof("Executing command `%s` in server `%s` ", command.Name(), data.ID)
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

func (c *AnimeModule) getAniUserIDFromMessage(m *discordgo.MessageCreate, s *discordgo.Session) (int, bool) {
	input := strings.SplitN(m.Content, " ", 2)

	// 3 Cases: no input, @mention, command input
	var userId int
	// Case: get saved id from database for that user
	if len(input) == 1 || len(m.Mentions) > 0 {
		var aniUserId int
		var errorMsg string
		// No user was mentioned, this means author is the target
		if len(input) == 1 {
			aniUserId = c.Bot.UserDataFromID(m.Author.ID).AniListData.UserId
			errorMsg = fmt.Sprintf("You have not yet linked an AniList user account, use `%saniset` to set one.", errorMsg)
		} else { // User was mentioned, make this user the target.
			aniUserId = c.Bot.UserDataFromID(m.Mentions[0].ID).AniListData.UserId
			errorMsg = fmt.Sprintf("This user has not yet linked an AniList account, use `%saniset` to set one.", errorMsg)
		}
		if aniUserId == 0 {
			// Case: Nothing set for this user
			s.ChannelMessageSend(m.ChannelID, errorMsg)
			return 0, true
		}
		userId = aniUserId

	} else { // Case: Message provided username, get username from query.
		userName := input[1]
		aniUser, err := queryUser(userName)
		if err != nil {
			// Case: No match found
			s.ChannelMessageSend(m.ChannelID, "Unable to find user with this name.")
			return 0, true
		}
		userId = aniUser.Id
	}
	return userId, false
}

func queryUser(userName string) (res anilistgo.User, err error) {
	query := "query ($search: String) { User (search: $search) " +
		"{ id name avatar {large} siteUrl } }"
	variables := struct {
		Search string `json:"search"`
	}{
		userName,
	}

	a, _ := anilistgo.New()
	res, err = a.User(query, variables)
	return
}
