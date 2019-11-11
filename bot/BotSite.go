package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sort"
	"time"
)

type HttpServer struct {
}

var res map[string]string

func (h *HttpServer) start() {
	var err error
	res = make(map[string]string)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var cur *mongo.Cursor
	cur, err = serverCollection.Find(ctx, bson.M{})

	if err != nil {
		log.Error(err)
	}

	for cur.Next(ctx) {
		var result ServerData
		err := cur.Decode(&result)
		if err != nil {
			log.Error(err)
		}
		HServer.updateServerCommands(result.ID, &result)
	}

	http.HandleFunc("/", h.handleServer())
	http.ListenAndServe(":80", nil)
}

func getUrlFromID(id string) string {
	b64 := base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(id))[:7])
	return b64[:5]
}

func (h *HttpServer) updateServerCommands(id string, data *ServerData) string {
	url := getUrlFromID(id)

	var keys []string

	for k := range data.CustomCommands {
		keys = append(keys, k)
	}

	var commandList string

	if len(keys) > 0 {
		commandList = "This is a list of all custom commands.\n"

		sort.Strings(keys)

		for _, v := range keys {
			commandList = fmt.Sprintf("%s\n %s - %s", commandList, v, data.CustomCommands[v].Content)
		}

	} else {
		commandList = fmt.Sprintf("There are no custom commands yet. Use `%saddcommand` to add your first command!", data.Key)
	}

	commandList = fmt.Sprintf("%s \n\nUse %scommandlist for a list of default commands.", commandList, data.Key)

	res[url] = commandList
	response := fmt.Sprintf("%s/%s", BotConfig.WebsiteUrl, url)
	return response
}

func (h *HttpServer) handleServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if commands, ok := res[r.URL.Path[1:]]; ok {
			fmt.Fprintf(w, commands)
		}

	}
}
