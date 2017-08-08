package main

import (
	"testing"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var basicUser = &discordgo.User{ID: "3333", Username: "Me"}

func TestParseCommand(t *testing.T) {
	m := &discordgo.MessageCreate{Message: &discordgo.Message{ID: "1111", ChannelID: "2222", Content: "!This is a message", Author: basicUser}}
	ret := parseCommand(m)

	if ret.Key != "!" {
		t.Errorf("Expected !, but got %s", ret.Key)
	}

	if ret.Command != "this" {
		t.Errorf("Expected this, but got %s", ret.Command)
	}

	actualContent := strings.Join(ret.Content[0:], " ")
	if actualContent != "is a message" {
		t.Errorf("Expected is a message, but got %s", actualContent)
	}

}