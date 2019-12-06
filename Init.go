package main

import (
	"flag"
	"github.com/gendonl/genbot/AnimeModule"
	"github.com/gendonl/genbot/Bot"
	"github.com/gendonl/genbot/CoreModule"
	"github.com/gendonl/genbot/MetaModule"
	"github.com/gendonl/genbot/ModerationModule"
	"github.com/gendonl/genbot/Server"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "./json/config.json", "<FILE>, where the configuration of the bot is located.")
	flag.Parse()
}

func main() {
	b := Bot.Bot{}
	modules := []Bot.Module{
		MetaModule.New(&b),
		CoreModule.New(&b),
		ModerationModule.New(&b),
		AnimeModule.New(&b),
	}
	server := Server.New(&b)
	b.InitBot(modules, server, configPath)
}
