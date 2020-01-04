package CoreModule

import "C"
import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"github.com/gendonl/saucenao-go"
	"net/url"
	"strconv"
	"strings"
)


func initSourceCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "source",
		description: "Gives the source of an image of either the url given as argument or the last image sent.",
		usage:       "`%ssource [image url]`",
		aliases:	 []string{"sauce"},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*CoreModule).sourceCommand,
	}
	return
}

func (c *CoreModule) sourceCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	input := strings.SplitN(m.Content, " ", 2)

	var imageUrl string
	if len(input) == 1 {
		amountOfMessages := 10
		res, err := c.Bot.FindLastMessageWithAttachOrEmbed(s, m, amountOfMessages)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unable to find an image to query in last %d messages.", amountOfMessages))
			return
		}
		imageUrl = res
	} else {
		imageUrl = input[1]
	}

	if !isValidUrl(imageUrl) {
		s.ChannelMessageSend(m.ChannelID, "Unable to parse command arguments to URL.")
		return
	}

	// Only query something that looks like an image.
	// We do not want to waste requests due to the API rate-limit.
	// FIXME: Use meta-data rather than url to determine whether something is an image.
	hasExtension := false
	for _, ext := range extensions {
		if strings.Contains(imageUrl, ext) {
			hasExtension = true
		}
	}

	if !hasExtension {
		_, _ = s.ChannelMessageSend(m.ChannelID, "URL has a non-supported file extension.")
		return
	}

	sauceClient := saucenao.New(c.Bot.Config.SauceNaoToken)

	sauceResult, err := sauceClient.FromURL(imageUrl)

	if err != nil || len(sauceResult.Data) == 0 {
		Log.Errorf("%s \n %s", err, sauceResult)
		_, _ = s.ChannelMessageSend(m.ChannelID, "Something went wrong while contacting the SauceNAO API." +
			"You could try sourcing your images manually at https://saucenao.com/")
		return
	}

	e := Bot.NewEmbed().
		SetAuthorFromUser(m.Author).
		SetColorFromUser(s, m.ChannelID, m.Author).
		SetThumbnail(imageUrl).
		SetFooter("This command is powered by SauceNAO.")

	similarity, err := strconv.ParseFloat(sauceResult.Data[0].Header.Similarity, 32)
	if err != nil || similarity < 80.0 {
		e.SetTitle("No results found")
		e.SetDescription("No images found with a confidence over 80.")
	} else {
		e.SetTitle(fmt.Sprintf("[Source] %s", sauceResult.Data[0].Data.Title)).
			SetURL(sauceResult.Data[0].Data.ExtUrls[0]).
			SetDescription(fmt.Sprintf("Result found with %s%% confidence", sauceResult.Data[0].Header.Similarity))
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}