package CoreModule

import (
	"fmt"
	"github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"strings"
	"time"
)


func initWeatherCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "weather",
		description: "Send the current weather for the given city.",
		usage:       "`%sweather <City>`",
		aliases:	 []string{},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*CoreModule).weatherCommand,
	}
	return
}

func (c *CoreModule) weatherCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	owm, err := openweathermap.NewCurrent("C", "EN", c.Bot.Config.OwmToken)

	if err != nil {
		Log.Error(err)
		return
	}

	input := strings.SplitN(m.Content, " ", 2)
	if len(input) > 2 {
		result := c.Bot.Usage(cmd, s, m, data)
		s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
		return
	}

	err = owm.CurrentByName(input[1])
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Either the OpenWeatherMap API is down or you provided an invalid location.")
		return
	}

	fahr := owm.Main.Temp*9/5 + 32

	// Convert timezone data from Seconds to hours
	GmtOffset := owm.Timezone / 60 / 60
	localTime := time.Now().UTC().Add(time.Duration(GmtOffset) * time.Hour).Format("3:04PM, Monday") // Local time

	// Generate the url of the weather icon
	iconUrl := "http://openweathermap.org/img/wn/" + owm.Weather[0].Icon + "@2x.png"

	// Convert degrees to wind direction.
	directionVal := int((owm.Wind.Deg / 22.5) + .5)
	directions := []string{"north", "north-northeast", "northeast", "east-northeast", "east", "east-southeast",
		"southeast", "south-southeast", "south", "south-southwest", "southwest", "west-southwest", "west", "west-northwest", "northwest", "north-northwest"}
	windDirection := directions[(directionVal % 16)]

	// Generate teh flag emoji
	flag := fmt.Sprintf(":flag_%s:", strings.ToLower(owm.Sys.Country))

	result := Bot.NewEmbed().
		SetAuthorFromUser(m.Author).
		SetColorFromUser(s, m.ChannelID, m.Author).
		SetThumbnail(iconUrl).
		SetTitle(fmt.Sprintf("Weather in **%s** at **%s** %s", owm.Name, localTime, flag)).
		AddField("Current Conditions:", fmt.Sprintf("**%s** at **%.1f°C** / **%.1f°F**",
			owm.Weather[0].Description, owm.Main.Temp, fahr)).
		AddInlineField("Humidity", fmt.Sprintf("%d%%", owm.Main.Humidity), true).
		AddInlineField("Wind", fmt.Sprintf("%.1f km/h from the %s ", owm.Wind.Speed*3.6, windDirection), true).
		SetFooter("Data provided by OpenWeatherMap", "http://f.gendo.moe/KlhvQJoD.png")

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, result.MessageEmbed)
}