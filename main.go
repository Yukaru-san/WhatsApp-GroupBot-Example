package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
	"golang.org/x/text/language"
)

// Settings - Change if you need more
type Settings struct {
	StdTranslationFrom language.Tag
	StdTranslationTo   language.Tag
}

var (
	conn *whatsapp.Conn
	sess whatsapp.Session

	// A List of groups that the default commands will work in
	standardGroups = []string{"Test-Group1", "Test-Group2"}
)

func main() {

	// Initialize Bot
	wabot.SetSessionFilePath("./data/storedSession.dat")
	wabot.SetUsersFilePath("./data/storedUsers.dat")
	wabot.SetQRFilePath("./data/qrCode.png")
	wabot.SetErrorTimeout(time.Minute)
	wabot.SetNicknameUseage(true, "Your nickname has been updated!") // TODO: Multi group support

	// Login
	var err error
	sess, conn, err = wabot.StartBot("github.com/rhymen/go-whatsapp", "go-whatsapp")

	if err != nil {
		fmt.Println("First login error:", err.Error())
	}

	// Add new savable user informations
	wabot.CreateNewSettingsOption(Settings{
		language.German,
		language.English,
	})

	// Try to load savedata and convert it into your own structure
	savedata, found := wabot.GetSaveData()
	if found {
		for i := 0; i < len(savedata.BotUsers); i++ {
			tmp := &Settings{}

			data, _ := json.Marshal(savedata.BotUsers[i].Settings)
			json.Unmarshal(data, tmp)

			savedata.BotUsers[i].Settings = *tmp
		}
		wabot.UseSaveData(savedata)
	}
	fmt.Println("Bot started!")

	// --------------------------[From a given standard list of groups]-----------------------------------------------

	// Handles a translation request
	wabot.AddGroupCommand("/t", standardGroups, func(message whatsapp.TextMessage) {
		go HandleTranslateRequest(message)
	})

	// Prints the user's settings
	wabot.AddGroupCommand("/printSettings", standardGroups, func(message whatsapp.TextMessage) {
		go HandleSettingsRequest(message)
	})

	// Picks a random entry from a given list
	wabot.AddGroupCommand("/pick", standardGroups, func(message whatsapp.TextMessage) {
		go HandlePickRequest(message)
	})

	// Prints the help message
	wabot.AddGroupCommand("/help", standardGroups, func(message whatsapp.TextMessage) {
		go wabot.WriteTextMessage(mainHelpMsg, message.Info.RemoteJid)
	})

	// Starts the dialog of playing a game
	wabot.AddGroupCommand("/play", standardGroups, func(message whatsapp.TextMessage) {
		go HandleGameRequest(message)
	})

	// --------------------------[Group Specific]-----------------------------------------------

	// Prints another help message when receiving the message in the given group
	wabot.AddGroupCommand("/ff20", []string{"{Some Group Name}"}, func(message whatsapp.TextMessage) {
		go wabot.WriteTextMessage(ff20HelpMsg, message.Info.RemoteJid)
	})

	// Sends a sticker upon writing "/m7" in the specified groups
	wabot.AddGroupCommand("/m7", []string{"{Some Group Name}", "{Another Group Name}"}, func(message whatsapp.TextMessage) {
		go HandleStickerRequest(message)
	})

	// Does something when receiving a sticker
	wabot.SetStickerHandler(func(message whatsapp.StickerMessage) {
		ReceivedStickerhandler(message)
	})

	// Don't let the program die here
	for {
		time.Sleep(time.Minute)
	}

}
