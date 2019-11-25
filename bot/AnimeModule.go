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
		Usage:       "`%sanime <anime name>",
		Permission:  discordgo.PermissionSendMessages,
		Execute:     animeInfoCommand,
	}

	cmd.AnimeCommands["anime"] = animeInfoCommand
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
	} else if cmd, ok := serverData.CustomCommands[msg.Command]; ok {
		log.Infof("Executing custom command `%s` in server `%s` ", cmd.Name, serverData.ID)

		s.ChannelMessageSend(m.ChannelID, cmd.Content)
	}
}

func animeInfoCommand(command Command, s *discordgo.Session, msg SentMessageData, data *ServerData) {
	if len(msg.Content) == 0 {
		result := createUsageInfo(command, msg, s, data)
		s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)
		return
	}

	animeName := strings.Join(msg.Content, " ")

	a, _ := anilistgo.New()

	res, err := a.Media(anilistgo.MediaVariables{
		SearchQuery: animeName,
	})

	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, "Unable to find anime with this name.")
		return
	}

	color, _ := strconv.ParseInt(strings.Replace(res.CoverImage.Color, "#", "", -1), 16, 32)

	description := strings.Split(res.Description, "<br>")[0]

	result := NewEmbed().
		SetAuthorFromUser(msg.Author).
		SetColor(int(color)).
		SetThumbnail(res.CoverImage.Medium).
		SetTitle(res.Title.English).
		SetURL(res.SiteUrl).
		SetDescription(description).
		AddInlineField("Episodes", strconv.Itoa(res.Duration), true).
		AddInlineField("Mean Score",  fmt.Sprintf("%.1f", float32(res.MeanScore)/10.0), true).
		SetFooter("This command is still under construction.")

	s.ChannelMessageSendEmbed(msg.ChannelID, result.MessageEmbed)
}