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

func (c *AnimeModule) getAniUserIDFromMessage(m *discordgo.MessageCreate) (userId int) {
	// Case: Check if we're able to get a user from the sent message.
	userData := Bot.UserDataFromMessage(m)
	if userData != nil {
		return userData.AniListData.UserId
	}

	// Message provided username, get username from query.
	input := strings.SplitN(m.Content, " ", 2)
	userName := input[1]
	aniUser, err := queryBaseUser(userName)
	if err != nil {
		// No match found
		return 0
	}
	userId = aniUser.Id
	return userId
}

func queryBaseUser(username string) (res anilistgo.User, err error) {
	query := "query ($search: String) { User (search: $search) " +
		"{ id name avatar {large} siteUrl } }"
	variables := struct {
		Search string `json:"search"`
	}{
		username,
	}

	a, _ := anilistgo.New()
	res, err = a.User(query, variables)
	return
}


func queryExtendedUser(username string, id int) (anilistgo.User, error) {
	var query string
	if id != 0 {
		query = "query ($id: Int) { User (id: $id) "
	} else {
		query = "query ($search: String) { User (search: $search) "
	}
	query = query + "{ id name avatar {large}  statistics {" +
		"anime {count episodesWatched} " +
		"manga {count chaptersRead}} siteUrl } }"
	variables := struct {
		Search string `json:"search"`
		Id int `json:"id"`
	}{
		username,
		id,
	}

	a, _ := anilistgo.New()
	res, err := a.User(query, variables)
	return res, err
}


func queryActivityData(aniUserID int) (res anilistgo.Activity, err error) {
	query := "query ($userid: Int) { Activity(userId: $userid, sort: ID_DESC) " +
		"{ ... on ListActivity { user { name } createdAt status progress media { type format" +
		" coverImage {large color} title { romaji native } " +
		"status episodes chapters siteUrl averageScore} } }  }  "

	variables2 := struct {
		Id int `json:"userid"`
	}{
		aniUserID,
	}

	a, _ := anilistgo.New()
	res, err = a.Activity(query, variables2)
	return
}
