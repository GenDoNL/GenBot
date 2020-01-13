package AnimeModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
	"time"
)

func initSetUserCommand() (cc AnimeCommand) {
	cc = AnimeCommand{
		name:        "aniset",
		description: "Links your discord account with an my anilist user.",
		usage:       "`%saniset <name/id>`",
		aliases:	 []string{},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*AnimeModule).setUserCommand,
	}
	return
}

func (c *AnimeModule) setUserCommand(cmd AnimeCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	if len(input) < 2 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	user := c.Bot.UserDataFromID(m.Author.ID)

	aniUser, err := queryUser(input[1])

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find user with this name.")
		return
	}

	user.AniListData = Bot.AniListUserData{
		UserId : aniUser.Id,
		LastUpdated: time.Now(),
	}

	user.WriteToDB()

	e := Bot.NewEmbed().SetAuthorFromUser(m.Author).
		SetTitle(fmt.Sprintf("[AniList] - %s", aniUser.Name)).
		SetImage(aniUser.Avatar.Large).
		SetURL(aniUser.SiteUrl)

	res := &discordgo.MessageSend{
		Embed: e.MessageEmbed,
		Content: "Set your AniList user to the following user:",
	}

	s.ChannelMessageSendComplex(m.ChannelID, res)
}