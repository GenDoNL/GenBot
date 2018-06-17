package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
	"strconv"
	"strings"
)

type OsuModule struct {
	test      string
	osuClient osuapi.Client
}

func (osu *OsuModule) setup() {
	osu.osuClient = *osuapi.NewClient(BotConfig.OsuToken)
	log.Info("osu module initialized.")
}

func (osu *OsuModule) execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	beatMapSetString := "https://osu.ppy.sh/s/"
	beatMapString := "https://osu.ppy.sh/b/"
	newSiteBeatMapString := "https://osu.ppy.sh/beatmapsets/"

	//check if message contains only one beatmap link
	beatmapSetID := strings.Count(m.Content, beatMapSetString)
	beatmapID := strings.Count(m.Content, beatMapString)
	newSiteBeatmapID := strings.Count(m.Content, newSiteBeatMapString)
	if beatmapSetID+beatmapID+newSiteBeatmapID != 1 {
		return
	}

	var opts osuapi.GetBeatmapsOpts
	var err error
	if beatmapSetID == 1 {
		opts, err = getBeatmapSetID(m.Content, beatMapSetString)
	}

	if beatmapID == 1 {
		opts, err = getBeatmapID(m.Content, beatMapString)
	}

	if newSiteBeatmapID == 1 {
		opts, err = getNewSiteBeatMapID(m.Content)
	}

	if err != nil {
		log.Errorf("Unable to parse to beatmap: \"%s\"", m.Content)
		return
	}

	var beatmap osuapi.Beatmap
	if beatmap, err = getBeatMap(osu.osuClient, opts); err != nil {
		log.Errorf("Something went wrong while trying to retrieve beatmap, %s", err)
		return
	}

	result := constructBeatmapMessage(beatmap)
	s.ChannelMessageSendEmbed(m.ChannelID, &result)
}

func getBeatmapSetID(content string, filter string) (osuapi.GetBeatmapsOpts, error) {
	//get beatmap id in the message
	s := content[(strings.Index(content, filter) + len(filter)):]
	id := strings.Split(s, " ")[0]
	beatMapSetID, err := strconv.Atoi(id)

	opts := osuapi.GetBeatmapsOpts{BeatmapSetID: beatMapSetID}
	return opts, err
}

func getBeatmapID(content string, filter string) (osuapi.GetBeatmapsOpts, error) {
	s := content[(strings.Index(content, filter) + len(filter)):]
	idString := strings.Split(s, " ")[0]
	idNoAnd := strings.Split(idString, "&")[0]
	id := strings.Split(idNoAnd, "?")[0]

	beatMapID, err := strconv.Atoi(id)

	opts := osuapi.GetBeatmapsOpts{BeatmapID: beatMapID}
	return opts, err
}

func getNewSiteBeatMapID(content string) (osuapi.GetBeatmapsOpts, error) {
	idString := strings.Split(content, "/")
	beatMapID, err := strconv.Atoi(idString[len(idString)-1])

	opts := osuapi.GetBeatmapsOpts{BeatmapID: beatMapID}
	return opts, err
}

func getBeatMap(osuClient osuapi.Client, opts osuapi.GetBeatmapsOpts) (osuapi.Beatmap, error) {
	//get beatMap info
	var beatMap osuapi.Beatmap
	beatMaps, err := osuClient.GetBeatmaps(opts)
	if err != nil {
		return beatMap, err
	}

	//check for empty list in case of no beatMaps found
	if len(beatMaps) == 0 {
		return beatMap, err
	}

	//get the highest difficulty in mapset
	beatMap = beatMaps[0]
	for i := 1; i < len(beatMaps); i++ {
		if beatMaps[i].DifficultyRating > beatMap.DifficultyRating {
			beatMap = beatMaps[i]
		}
	}

	return beatMap, nil
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
