package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/anilist-go"
	"strconv"
	"strings"
	"fmt"
)

type AnimeModule struct {
	AnimeCommands map[string]Command
	HelpString		string
}

func (cmd *AnimeModule) setup() {

	cmd.AnimeCommands = map[string]Command{}

	animeInfoCommand := Command{
		Name:        "anime",
		Description: "Sends info on the given anime.",
		Usage:       "`%sanime <anime name>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     animeInfoCommand,
	}
	cmd.AnimeCommands["anime"] = animeInfoCommand

	mangaInfoCommand := Command{
		Name:        "manga",
		Description: "Sends info on the given manga.",
		Usage:       "`%smanga <manga name>`",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     mangaInfoCommand,
	}

	cmd.AnimeCommands["manga"] = mangaInfoCommand
}

func (cmd *AnimeModule) retrieveCommands() map[string]Command {
	return cmd.AnimeCommands
}

func (cmd *AnimeModule) retrieveHelp() (moduleName string, info string) {
	moduleName = "Anime Module"
	info = commandsToHelp(&cmd.HelpString, cmd.AnimeCommands)
	return
}

func (cmd *AnimeModule) execute(s *discordgo.Session, m *discordgo.MessageCreate, msg SentMessageData, serverData *ServerData) {
	if serverData.Key != msg.Key {
		return
	}

	if command, ok := cmd.AnimeCommands[msg.Command]; ok {
		isCommander, ok := serverData.Commanders[m.Author.ID]
		perm, _ := s.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
		if perm&command.Permission == command.Permission || (ok && isCommander) {
			log.Infof("Executing command `%s` in server `%s` ", command.Name, serverData.ID)
			command.Execute(command, s, msg, serverData)
		} else {
			log.Infof("Use of %s command denied for permission level %d", command.Name, perm)
		}
	}
}

func mediaInfoCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(msg.Content) == 0 {
		result := createUsageInfo(command, msg, s, data)
		s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)
		return
	}
	mediaType := strings.ToUpper(command.Name)
	name := strings.Join(msg.Content, " ")

	res, err := queryBasicMediaInfo(name, mediaType, s, msg)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, "Unable to find anime with this name.")
		return
	}

	color, _ := strconv.ParseInt(strings.Replace(res.CoverImage.Color, "#", "", -1), 16, 32)

	score := fmt.Sprintf("%.1f", float32(res.MeanScore)/10.0)
	status := strings.Title(strings.ToLower(res.Status))
	status = strings.Replace(status, "_", " ", -1)

	title := parseTitle(res)
	description := parseDescription(res)
	episodeChapters := chaptersOrEpisodes(mediaType, res)

	e := NewEmbed().
		SetColor(int(color)).
		SetThumbnail(res.CoverImage.Large).
		SetTitle(title).
		SetURL(res.SiteUrl).
		SetDescription(description).
		SetFooter(fmt.Sprintf("Score: %s    Status: %s    %s", score, status, episodeChapters))

	s.ChannelMessageSendEmbed(msg.ChannelID, e.MessageEmbed)
}

func queryBasicMediaInfo(name string, mediaType string, s *discordgo.Session, msg SentMessageData) (res anilistgo.Media, err error) {
	query := "query ($search: String, $type: MediaType) { Media (search: $search, type: $type) " +
		"{ id description(asHtml: false) coverImage {large color} title { romaji native } status episodes chapters siteUrl meanScore} }"
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
			episodeChapters = fmt.Sprintf("Episodes: unknown")
		} else {
			episodeChapters = fmt.Sprintf("Episodes: %d", res.Episodes)
		}
	} else {
		if res.Chapters == 0 {
			episodeChapters = fmt.Sprint("Chapters: unknown")
		} else {
			episodeChapters = fmt.Sprintf("Chapters: %d", res.Chapters)
		}
	}
	return episodeChapters
}

func animeInfoCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	mediaInfoCommand(command, s, msg, data)

}

func mangaInfoCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	mediaInfoCommand(command, s, msg, data)
}