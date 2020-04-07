package RSSModule

import (
	"context"
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/gendonl/genbot/Bot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type RSSPoll struct {
	Url string `bson:"url"`
	LastUpdated time.Time `bson:"last-updated"`
	Channels []string `bson:"channels"`
}

func (c *RSSModule) initPolling() {
	// Whatever happens, do NOT go down.
	defer func() {
		if r := recover(); r != nil {
			Log.Criticalf("Bot panicked while polling: ", r)
		}
	}()

	time.Sleep(5 * time.Second)

	// Start polling
	for {
		c.poll()

		// Lazy update once every 60 second
		// If functionality is used more often, this has to be updated to reduce latency
		time.Sleep(60 * time.Second)
	}
}

func (c *RSSModule) poll() {
	if Bot.RSSCollection == nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{}
	cursor, err := Bot.RSSCollection.Find(ctx, filter)
	if err != nil {
		Log.Error(err)
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result RSSPoll
		err := cursor.Decode(&result)
		if err != nil {
			Log.Error(err)
			continue
		}

		c.checkUpdates(result)
	}
}

func (c *RSSModule) checkUpdates(poll RSSPoll) {
	feed, err := rss.Fetch(poll.Url)

	var toBeSend []*rss.Item
	if err != nil || len(feed.Items) == 0  {
		return
	}

	for _, newItem := range feed.Items {
		if poll.LastUpdated.After(newItem.Date) {
			break
		}
		toBeSend = append(toBeSend, newItem)
	}

	err = c.updateTime(poll.Url)
	if err != nil {
		Log.Error(err)
		return
	}
	if len(toBeSend) > 0 {
		c.sendUpdate(toBeSend, poll)
	}
}

func (c *RSSModule) sendUpdate(items []*rss.Item, poll RSSPoll) {
	for i := len(items)-1; i >= 0; i-- {
		item := items[i]
		for _, channel := range poll.Channels {
			embed := createEmbed(item)
			c.Bot.Session.ChannelMessageSendEmbed(channel, embed.MessageEmbed)
		}
	}
}

func createEmbed(item *rss.Item) (e *Bot.Embed){
	footer := fmt.Sprintf(item.Date.Format("Monday, 02-Jan, 3:04PM"))
	e = Bot.NewEmbed().SetTitle(item.Title).SetURL(item.Link).SetDescription (item.Summary).SetFooter(footer)
	for _, enc := range item.Enclosures {
		if strings.HasPrefix(enc.Type, "image") {
			e.SetThumbnail(enc.URL)
			break
		}
	}
	return
}

func (c *RSSModule) updateTime(url string) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	_, err = Bot.RSSCollection.UpdateOne(
		ctx,
		bson.D{
			{"url", url },
		},
		bson.D{
			{
				"$set", bson.M{
				"last-updated": time.Now(),
				},
			},
		},
		options.Update().SetUpsert(true))

	return
}
