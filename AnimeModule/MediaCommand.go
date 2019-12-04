package AnimeModule

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	anilistgo "github.com/gendonl/anilist-go"
	"github.com/gendonl/genbot/Bot"
	"strconv"
	"strings"
)

func initAnimeCommand() (cc AnimeCommand) {
	cc = AnimeCommand{
		name:        "anime",
		description: "Returns information on the provided anime",
		usage:       "`%sanime <anime name>`",
		aliases:	 []string{"a"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*AnimeModule).mediaCommand,
	}
	return
}

func initMangaCommand() (cc AnimeCommand) {
	cc = AnimeCommand{
		name:        "manga",
		description: "Returns information on the provided manga",
		usage:       "`%smanga <manga name>`",
		aliases:	 []string{"m"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*AnimeModule).mediaCommand,
	}
	return
}

func (c *AnimeModule) mediaCommand(cmd AnimeCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	if len(input) < 2 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	// Command name is either "manga" or "anime", therefore we can safely use it as MediaType in the query.
	mediaType := strings.ToUpper(cmd.Name())
	contentName := input[1]

	res, err := queryBasicMediaInfo(contentName, mediaType)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to find %s with this name.", cmd.Name()))
		return
	}

	color, _ := strconv.ParseInt(strings.Replace(res.CoverImage.Color, "#", "", -1), 16, 32)

	score := fmt.Sprintf("%.1f", float32(res.AverageScore)/10.0)
	status := strings.Title(strings.ToLower(res.Status))
	status = strings.Replace(status, "_", " ", -1)

	title := parseTitle(res)
	description := parseDescription(res)
	episodeChapters := chaptersOrEpisodes(mediaType, res)

	e := Bot.NewEmbed().
		SetColor(int(color)).
		SetThumbnail(res.CoverImage.Large).
		SetTitle(title).
		SetURL(res.SiteUrl).
		SetDescription(description).
		SetFooter(fmt.Sprintf("Score: %s    Status: %s    %s    Type: %s", score, status, episodeChapters, strings.Title(cmd.Name())))

	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)

}

func queryBasicMediaInfo(name string, mediaType string) (res anilistgo.Media, err error) {
	query := "query ($search: String, $type: MediaType) { Media (search: $search, type: $type) " +
		"{ id description(asHtml: false) coverImage {large color} title { romaji native } status episodes chapters siteUrl averageScore} }"
	variables := struct {
		Search string `json:"search"`
		Type   string `json:"type"`
	}{
		name,
		mediaType,
	}

	a, _ := anilistgo.New()
	res, err = a.Media(query, variables)

	return
}

func parseTitle(res anilistgo.Media) string {
	var title string
	if res.Title.Romaji != "" {
		title = res.Title.Romaji
	} else {
		title = res.Title.Native
	}
	return title
}

func parseDescription(res anilistgo.Media) string {
	description := strings.Split(res.Description, "<br>")[0]
	description = strings.Replace(description, "<i>", "*", -1)
	description = strings.Replace(description, "</i>", "*", -1)
	return description
}

func chaptersOrEpisodes(mediaType string, res anilistgo.Media) string {
	var episodeChapters string
	if mediaType == "ANIME" {
		if res.Episodes == 0 {
			episodeChapters = fmt.Sprintf("Episodes: Unknown")
		} else {
			episodeChapters = fmt.Sprintf("Episodes: %d", res.Episodes)
		}
	} else {
		if res.Chapters == 0 {
			episodeChapters = fmt.Sprint("Chapters: Unknown")
		} else {
			episodeChapters = fmt.Sprintf("Chapters: %d", res.Chapters)
		}
	}
	return episodeChapters
}