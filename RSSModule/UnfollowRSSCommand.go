package RSSModule

import (
	"context"
	"github.com/SlyMarbo/rss"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func initUnfollowRSSCommand() (cc RSSCommand) {
	cc = RSSCommand{
		name:        "unfollow",
		description: "Unfollows an RSS feed.",
		usage:       "`%sfollow <feed_url>`",
		aliases:     []string{"rmfeed", "delfeed"},
		permission:  discordgo.PermissionManageServer,
		execute:     (*RSSModule).UnfollowRSS,
	}
	return
}

func (c *RSSModule) UnfollowRSS(cmd RSSCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	input := strings.SplitN(m.Content, " ", 2)

	if len(input) > 2 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	feed, err := rss.Fetch(input[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while retrieving feed.")
		return
	}

	err = c.rmRSSFollow(feed.UpdateURL, m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while retrieving feed.")
		return
	}

	e := Bot.NewEmbed().SetTitle(feed.Title).SetURL(feed.UpdateURL).SetDescription("You unfollowed this RSS.")
	if feed.Image != nil {
		e.SetThumbnail(feed.Image.URL)
	}
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (c *RSSModule) rmRSSFollow(url string, channelid string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	_, err = Bot.RSSCollection.UpdateOne(
		ctx,
		bson.D{
			{"url", url },
		},
		bson.D{
			{
				"$pull", bson.M{
				"channels": channelid,
			},
			},
		},
		options.Update().SetUpsert(false))

	return
}