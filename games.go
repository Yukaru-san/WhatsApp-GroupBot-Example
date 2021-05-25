package main

import (
	"strings"

	"github.com/Rhymen/go-whatsapp"
	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
)

// HandleGameRequest is used for /play commands
func HandleGameRequest(message whatsapp.TextMessage) {

	// Input check
	if len(message.Text) < 7 {
		wabot.WriteTextMessage("Please tell me your desired game in the format /play {game}.\n\nPossible games are:\nhangman", message.Info.RemoteJid)
		return
	}

	// Request playing hangman
	if strings.HasPrefix(strings.ToLower(message.Text[6:]), "hangman") {
		if hangmanIsRunning {
			wabot.WriteTextMessage("Sorry!\nThis game is already in progress.", message.Info.RemoteJid)
		} else {
			InitiateHangman(message)
		}
	}
}
