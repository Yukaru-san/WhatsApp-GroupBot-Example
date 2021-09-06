package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/Rhymen/go-whatsapp"
	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
)

// HandleStickerRequest replies with requested Stickers
func HandleStickerRequest(message whatsapp.TextMessage) {
	if strings.HasPrefix(message.Text, "/m7") {
		img, _ := os.Open("stickers/{somesticker}.webp")

		wabot.SendStickerMessage(img, message.Info.RemoteJid)
	}
}

// ReceivedStickerhandler can react to sent stickers and echo them back (I used it in a chat with a specific person)
func ReceivedStickerhandler(message whatsapp.StickerMessage) {
	if message.Info.RemoteJid == wabot.NameToJid("{Some contact's name}") {
		HandleStickerEcho(message)
	}
}

// HandleStickerEcho echoes the same sticker back to the user
func HandleStickerEcho(message whatsapp.StickerMessage) {
	// Download the image and store it
	c, _ := message.Download()
	ioutil.WriteFile("sticker.webp", c, 0600)

	// Open it to have a *os.File
	img, _ := os.Open("sticker.webp")

	// Send it
	wabot.SendStickerMessage(img, message.Info.RemoteJid)
}
