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
		aliases:	 []string{"ar", "arecent"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*AnimeModule).AniRecentCommand,
	}
	return
}


func (c *AnimeModule) AniRecentCommand(cmd AnimeCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {

	// Get the AniID of the targeted user, if notFound it means that no match was found
	userId := c.getAniUserIDFromMessage(m)
	if userId == 0 {
		s.ChannelMessageSend(m.ChannelID, "Unable to find AniList account with this name.")
		return
	}

	recentStatus, err := queryActivityData(userId)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find AniList account with this name.")
		Log.Error(err)
		return
	}

	e := createRecentActivityEmbed(recentStatus)

	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)

}

func createRecentActivityEmbed(recentStatus anilistgo.Activity) *Bot.Embed {
	media := recentStatus.Media

	color, _ := strconv.ParseInt(strings.Replace(media.CoverImage.Color, "#", "", -1), 16, 32)

	score := fmt.Sprintf("%.1f", float32(media.AverageScore)/10.0)
	status := strings.Title(strings.ToLower(media.Status))
	status = strings.Replace(status, "_", " ", -1)

	title := parseTitle(media)

	var description string
	if recentStatus.Status == "completed" {
		description = fmt.Sprintf("%s recently completed.", recentStatus.User.Name)
	} else if recentStatus.Status != "" {
		// TODO: Remove "of" on dropped anime
		description = fmt.Sprintf("%s recently %s %s", recentStatus.User.Name, recentStatus.Status, recentStatus.Progress)
	}

	episodeChapters := chaptersOrEpisodes(media)

	e := Bot.NewEmbed().
		SetColor(int(color)).
		SetImage(media.CoverImage.Large).
		SetTitle(title).
		SetURL(media.SiteUrl).
		SetDescription(description).
		SetFooter(fmt.Sprintf("Score: %s    Status: %s    Type: %s",
			score, episodeChapters, strings.Title(strings.ToLower(media.Format))))

	return e
}

