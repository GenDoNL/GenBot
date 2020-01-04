package main

import (
	"flag"
	"github.com/gendonl/genbot/AnimeModule"
	"github.com/gendonl/genbot/Bot"
	"github.com/gendonl/genbot/CoreModule"
	"github.com/gendonl/genbot/MetaModule"
	"github.com/gendonl/genbot/ModerationModule"
	"github.com/gendonl/genbot/Server"
	"github.com/op/go-logging"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "./json/config.json", "<FILE>, where the configuration of the bot is located.")
	flag.Parse()
}

func initLogger() (Log *logging.Logger){
	Log = logging.MustGetLogger("GenBot")
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	backEnd := logging.NewLogBackend(os.Stderr, "", 0)
	backEndFormatter := logging.NewBackendFormatter(backEnd, format)

	logging.SetBackend(backEndFormatter)

	Log.Info("Logger initialized.")
	return
}

func main() {
	log := initLogger()
	b := Bot.New(log)
	modules := []Bot.Module{
		MetaModule.New(b, log),
		CoreModule.New(b, log),
		ModerationModule.New(b, log),
		AnimeModule.New(b, log),
	}
	server := Server.New(b, log)
	b.InitBot(modules, server, configPath)
}
