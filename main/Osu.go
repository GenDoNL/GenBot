package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
	"strconv"
	"strings"
)
var (
	beatmapStringLength = 21
)

func getBeatmapSetID(content string) (int, error) {
	//get beatmap id in the message
	id, err := strconv.Atoi(strings.Split(string(content[strings.Index(content, "https://osu.ppy.sh/s/")+beatmapStringLength:]), " ")[0])
	if err != nil {
		fmt.Println("An error occurred while getting beatmap id, ", err)
		return -1, err
	}
	return id, nil
}

func getBeatmapID(content string) (int, error) {

		//get beatmap id in the message
	str := strings.Split(string(content[strings.Index(content, "https://osu.ppy.sh/b/")+beatmapStringLength:]), " ")[0]
	str = strings.Split(str, "?")[0]
	id, err := strconv.Atoi(strings.Split(str, "&")[0])
	if err != nil {
		fmt.Println("An error occurred while getting beatmap id, ", err)
		return -1, err
	}
	return id, nil
}

func constructBeatmapMessage(beatmap osuapi.Beatmap) discordgo.MessageEmbed {
	//create the embed message to send
	return discordgo.MessageEmbed{
		URL:         "https://osu.ppy.sh/s/" + strconv.Itoa(beatmap.BeatmapSetID),
		Title:       beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]",
		Description: "**Mode:** " + beatmap.Mode.String() + " | **Length:** " + parseTime(beatmap.TotalLength) + "\n**Star rating:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64) + " | **BPM:** " + strconv.FormatFloat(beatmap.BPM, 'f', 2, 64) + "\n**OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 2, 64) + " | **CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 2, 64) + "\n**AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 2, 64) + " | **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 2, 64),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
		},
		Color: 16763135, //this should be pink
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Mapped by " + beatmap.Creator + " | Status: " + beatmap.Approved.String(),
		},
	}
}

func getBeatmap(opts osuapi.GetBeatmapsOpts) (osuapi.Beatmap, error) {
	//get beatmap info
	var beatmap osuapi.Beatmap
	beatmaps, err := osuClient.GetBeatmaps(opts)
	if err != nil {
		return beatmap, err
	}

	//check for empty list in case of no beatmaps found
	if len(beatmaps) == 0 {
		return beatmap, err
	}

	//get the highest difficulty in mapset
	beatmap = beatmaps[0]
	for i := 1; i < len(beatmaps); i++ {
		if beatmaps[i].DifficultyRating > beatmap.DifficultyRating {
			beatmap = beatmaps[i]
		}
	}

	return beatmap, nil
}

func checkBeatmapLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	//check if message contains only one beatmap link
	beatmapSetID := strings.Count(m.Content, "https://osu.ppy.sh/s/")
	beatmapID := strings.Count(m.Content, "https://osu.ppy.sh/b/")
	if beatmapSetID + beatmapID != 1 {
		return
	}

	id := -1
	var err error
	var opts osuapi.GetBeatmapsOpts


	if beatmapSetID == 1 {
		id, err = getBeatmapSetID(m.Content)
		if err != nil {
			return
		}
		opts = osuapi.GetBeatmapsOpts{BeatmapSetID: id}

	} else if beatmapID == 1 {
		id, err = getBeatmapID(m.Content)
		if err != nil {
			return
		}
		opts = osuapi.GetBeatmapsOpts{BeatmapID: id}
	}

	var beatmap osuapi.Beatmap
	if beatmap, err = getBeatmap(opts); err != nil {
		fmt.Println("Something went wrong while trying to retrieve beatmap, ", err)
		return
	}

	message := constructBeatmapMessage(beatmap)

	//send the message. finally.
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &message)
	if err != nil {
		fmt.Println("An error occurred while sending embed message, ", err)
	}
}

//get time in format MM:SS
func parseTime(s int) string {
	if s < 10 {
		return "00:0" + strconv.Itoa(s)
	}
	if s < 61 {
		return "00:" + strconv.Itoa(s)
	}
	m := int(s / 60)
	s = s % 60
	if s < 10 {
		return strconv.Itoa(m) + ":0" + strconv.Itoa(s)
	}

	return strconv.Itoa(m) + ":" + strconv.Itoa(s)
}
