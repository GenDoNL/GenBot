package AnimeModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	anilistgo "github.com/gendonl/anilist-go"
	"github.com/gendonl/genbot/Bot"
	"strings"
)



func initAniUserInfoCommand() (cc AnimeCommand) {
	cc = AnimeCommand{
		name:        "aniuser",
		description: "Returns AniList information on the provided user",
		usage:       "`%saniuser <name/id>`",
		aliases:	 []string{"m"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*AnimeModule).AniUserInfoCommand,
	}
	return
}

func (c *AnimeModule) AniUserInfoCommand(cmd AnimeCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	if len(input) < 2 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	userName := input[1]

	query := "query ($search: String) { User (search: $search) " +
		"{ id name avatar {large}  statistics {anime {count episodesWatched} manga {count chaptersRead}} siteUrl } }"
	variables := struct {
		Search string `json:"search"`
	}{
		userName,
	}

	a, _ := anilistgo.New()
	res, err := a.User(query, variables)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find user with this name.")
		return
	}

	description := fmt.Sprintf("**Anime Watched:** %d\n **Episodes Watched:** %d\n\n **Manga Read:** %d\n **Chapters Read:** %d\n",
		res.Statistics.Anime.Count, res.Statistics.Anime.EpisodesWatched,
		res.Statistics.Manga.Count, res.Statistics.Manga.ChaptersRead)

	recentStatus, err := c.getActivityData(res)

	if err == nil {
		description = fmt.Sprintf("%s\n%s", description, recentStatus)
	}

	e := Bot.NewEmbed().
		SetColorFromUser(s, m.ChannelID, m.Author).
		SetTitle(res.Name).
		SetDescription(res.About).
		SetThumbnail(res.Avatar.Large).
		SetURL(res.SiteUrl).
		SetDescription(description)

	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)

}

func (c *AnimeModule) getActivityData(usr anilistgo.User) (res string, err error) {
	query := "query ($userid: Int) { Activity(userId: $userid, sort: ID_DESC) " +
		"{ ... on ListActivity { createdAt status progress media { type title { romaji } }  }  } } "
	variables2 := struct {
		Id int `json:"userid"`
	}{
		usr.Id,
	}

	a, _ := anilistgo.New()
	activity, err := a.Activity(query, variables2)
	if err != nil {
		return "", err
	}

	if activity.Status == "completed" {
		res = fmt.Sprintf("Recently completed the %s **%s**",
			strings.ToLower(activity.Media.Type), activity.Media.Title.Romaji)
	} else if activity.Status != "" {
		// TODO: Remove "of" on dropped anime
		res = fmt.Sprintf("Recently %s %s of **%s**", activity.Status, activity.Progress,
			activity.Media.Title.Romaji)
	}
	return
}
