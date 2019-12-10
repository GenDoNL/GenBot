package AnimeModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	anilistgo "github.com/gendonl/anilist-go"
	"github.com/gendonl/genbot/Bot"
	"strconv"
	"strings"
)



func initAniRecentCommand() (cc AnimeCommand) {
	cc = AnimeCommand{
		name:        "anirecent",
		description: "Returns info on the most recent activity on the users AniList",
		usage:       "`%sanirecent <name/id>`",
		aliases:	 []string{"ar"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*AnimeModule).AniRecentCommand,
	}
	return
}

func (c *AnimeModule) AniRecentCommand(cmd AnimeCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	if len(input) < 2 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	userName := input[1]

	query := "query ($search: String) { User (search: $search) " +
		"{ id } }"
	variables := struct {
		Search string `json:"search"`
	}{
		userName,
	}

	a, _ := anilistgo.New()
	res, err := a.User(query, variables)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find user with this name.")
		c.Bot.Log.Error(err)
		return
	}


	recentStatus, err := getActivityData(res)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find user with this name.")
		c.Bot.Log.Error(err)
		return
	}
	media := recentStatus.Media

	color, _ := strconv.ParseInt(strings.Replace(media.CoverImage.Color, "#", "", -1), 16, 32)

	score := fmt.Sprintf("%.1f", float32(media.AverageScore)/10.0)
	status := strings.Title(strings.ToLower(media.Status))
	status = strings.Replace(status, "_", " ", -1)


	title := parseTitle(media)

	var description string
	if recentStatus.Status == "completed" {
		description = fmt.Sprintf("Recently completed.")
	} else if recentStatus.Status != "" {
		// TODO: Remove "of" on dropped anime
		description = fmt.Sprintf("Recently %s %s", recentStatus.Status, recentStatus.Progress)
	}

	episodeChapters := chaptersOrEpisodes(media)

	e := Bot.NewEmbed().
		SetColor(int(color)).
		SetImage(media.CoverImage.Large).
		SetTitle(title).
		SetURL(media.SiteUrl).
		SetDescription(description).
		SetFooter(fmt.Sprintf("Score: %s    Status: %s    Type: %s",
			score, episodeChapters, strings.Title(strings.ToLower(media.Type))))

	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)

}

func getActivityData(usr anilistgo.User) (res anilistgo.Activity, err error) {
	query := "query ($userid: Int) { Activity(userId: $userid, sort: ID_DESC) " +
		"{ ... on ListActivity { createdAt status progress media { type " +
		" coverImage {large color} title { romaji native } " +
		"status episodes chapters siteUrl averageScore} } }  }  "

	variables2 := struct {
		Id int `json:"userid"`
	}{
		usr.Id,
	}

	a, _ := anilistgo.New()
	res, err = a.Activity(query, variables2)
	return
}
