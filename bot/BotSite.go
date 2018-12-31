package main

import (
	"net/http"
	"fmt"
	"encoding/base64"
	"crypto/sha1"
)

type HttpServer struct {
	
}

var res map[string]string

func (h *HttpServer) start() {
	res = make(map[string]string)
	http.HandleFunc("/", h.handleServer())
	http.ListenAndServe(":80", nil)
}

func getUrlFromID(id string) string {
	b64 := base64.StdEncoding.EncodeToString(sha1.New().Sum([]byte(id))[:7])
	return b64[:5]
}

func (h *HttpServer) updateServerCommands(id string, result string) string {
	url := getUrlFromID(id)
	res[url] = result
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