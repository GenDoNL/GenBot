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

	// Fill the map will all known servers and their hashes
	for cur.Next(ctx) {
		var result ServerData
		err := cur.Decode(&result)
		if err != nil {
			log.Error(err)
		}

		url := h.getPathFromID(result.ID)
		res[url] = result.ID
	}

	http.HandleFunc("/", h.handleServer())
	http.ListenAndServe(":80", nil)
}

func (h *HttpServer) getPathFromID(guildID string) string {
	b64 := base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(guildID))[:7])
	return b64[:5]
}

func (h *HttpServer) getUrlFromID(guildID string) string {
	url := h.getPathFromID(guildID)
	return fmt.Sprintf("%s/%s", BotConfig.WebsiteUrl, url)
}

// Convert the data of a server into a string
func (h *HttpServer) getSite(guildID string) string {
	var keys []string

	data := getServerData(guildID)

	for k := range data.CustomCommands {
		keys = append(keys, k)
	}

	var commandList string

	if len(keys) > 0 {
		sort.Strings(keys)

		for _, v := range keys {
			commandList = fmt.Sprintf("%s\n %s - %s", commandList, v, data.CustomCommands[v].Content)
		}

	} else {
		commandList = fmt.Sprintf("There are no custom commands yet. Use `%saddcommand` to add your first command!", data.Key)
	}

	return commandList
}

func (h *HttpServer) handleServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if id, ok := res[r.URL.Path[1:]]; ok {
			fmt.Fprintf(w, h.getSite(id))
		}

	}
}
