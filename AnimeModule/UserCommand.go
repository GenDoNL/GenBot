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
		aliases:	 []string{"anilist"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*AnimeModule).AniUserInfoCommand,
	}
	return
}

func (c *AnimeModule) AniUserInfoCommand(cmd AnimeCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	userData := Bot.UserDataFromMessage(m)
	var id int
	var username string
	if userData != nil {
		id = userData.AniListData.UserId
	} else {
		username = input[1]
	}

	res, err := queryExtendedUser(username, id)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find user with this name.")
		return
	}

	description := fmt.Sprintf("**Anime Watched:** %d\n **Episodes Watched:** %d\n\n **Manga Read:** %d\n **Chapters Read:** %d\n",
		res.Statistics.Anime.Count, res.Statistics.Anime.EpisodesWatched,
		res.Statistics.Manga.Count, res.Statistics.Manga.ChaptersRead)

	recentActivity, err := queryActivityData(res.Id)
	recentStatus, err := parseActivityString(recentActivity)

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



func parseActivityString(activity anilistgo.Activity) (res string, err error) {
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
