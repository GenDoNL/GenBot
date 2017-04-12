package main

import (
	"github.com/thehowl/go-osuapi"
	"github.com/bwmarrin/discordgo"
	"fmt"
	"strings"
	"strconv"
)

func checkBeatmapLink(s *discordgo.Session, m *discordgo.MessageCreate) {

	//check if message contains only one beatmap link
	BeatmapSet := strings.Count(m.Content, "https://osu.ppy.sh/s/")
	Beatmap := strings.Count(m.Content, "https://osu.ppy.sh/b/")
	if BeatmapSet + Beatmap != 1 {
		return
	}

	id := 0

	//check if message contains beatmap set
	if BeatmapSet == 1 {

		//get beatmap id in the message
		tmp_id, err := strconv.Atoi(strings.Split(string(m.Content[strings.Index(m.Content,"https://osu.ppy.sh/s/")+21:])," ")[0])
		if err != nil {
			fmt.Println("An error ocurred while getting beatmap id, ",err)
			return
		}
		id = tmp_id
	}

	//check if message contains specific beatmap difficulty
	if Beatmap == 1 {

		//get beatmap id in the message
		tmp_id, err := strconv.Atoi(strings.Split(strings.Split(string(m.Content[strings.Index(m.Content,"https://osu.ppy.sh/b/")+21:])," ")[0],"?")[0])
		if err != nil {
			fmt.Println("An error ocurred while getting beatmap id, ",err)
			return
		}
		id = tmp_id
	}

	//return if there's no beatmap link
	if id == 0 {
		return
	}

	var opts osuapi.GetBeatmapsOpts

	if BeatmapSet == 1 {
		opts = osuapi.GetBeatmapsOpts{BeatmapSetID: id}
	} else {
		opts = osuapi.GetBeatmapsOpts{BeatmapID: id}
	}

	//get beatmap info
	beatmaps, err := osu_client.GetBeatmaps(opts)
	if err != nil {
		fmt.Println("An error ocurred while fecthing beatmap, ",err)
		return
	}

	//check for empty list in case of no beatmaps found
	if len(beatmaps) == 0 {
		return
	}

	//get the highest difficulty in mapset
	beatmap := beatmaps[0]
	for i := 1; i < len(beatmaps); i++ {
		if beatmaps[i].DifficultyRating > beatmap.DifficultyRating {
			beatmap = beatmaps[i]
		}
	}

	//create the embed message to send
	message := discordgo.MessageEmbed{
		URL: "https://osu.ppy.sh/s/"+strconv.Itoa(beatmap.BeatmapSetID),
		Title: beatmap.Artist+" - "+beatmap.Title+" ["+beatmap.DiffName+"]",
		Description: "**Mode:** "+beatmap.Mode.String()+" | **Length:** "+parseTime(beatmap.TotalLength)+"\n**Star rating:** "+strconv.FormatFloat(beatmap.DifficultyRating,'f',2,64)+" | **BPM:** "+strconv.FormatFloat(beatmap.BPM,'f',2,64)+"\n**OD:** "+strconv.FormatFloat(beatmap.OverallDifficulty,'f',2,64)+" | **CS:** "+strconv.FormatFloat(beatmap.CircleSize,'f',2,64)+"\n**AR:** "+strconv.FormatFloat(beatmap.ApproachRate,'f',2,64)+" | **HP:** "+strconv.FormatFloat(beatmap.HPDrain,'f',2,64),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/"+strconv.Itoa(beatmap.BeatmapSetID)+"l.jpg",
		},
		Color: 16763135, //this should be pink
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Mapped by "+beatmap.Creator+" | Status: "+beatmap.Approved.String(),
		},
	}

	//send the message. finally.
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &message)
	if err != nil {
		fmt.Println("An error ocurred while sending embed message, ",err)
	}
}

//get time in format MM:SS
func parseTime(s int) string {

	if s < 10 {
		return "00:0"+strconv.Itoa(s)
	}
	if s < 61 {
		return "00:"+strconv.Itoa(s)
	}
	m := int(s/60)
	s = s - m * 60
	if s < 10 {
		return strconv.Itoa(m)+":0"+strconv.Itoa(s)
	}
	return strconv.Itoa(m)+":"+strconv.Itoa(s)
}