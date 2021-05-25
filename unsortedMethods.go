package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
)

var (
	mainHelpMsg = "Possible commands:\n" +
		"/nick\n   -   changes your nickname\n" +
		"/t\n   -   translate using default langs\n" +
		"/t-{from}-{to}\n   -   translate using the given langs\n" +
		"      e.g:   /t-de-es\n" +
		"/help\n   -   well.. it prints this msg\n" +
		"/tf {language}\n   -   set your default {from} lang\n" +
		"/tt {language}\n   -   set your default {to} lang\n" +
		"/ts-{from}-{to}\n   -   set your default {from} and {to}\n" +
		"/printSettings\n   -   prints your current settings\n" +
		"/al {day}\n   -    airing anime from today+{0-6}\n" +
		"/pick\n   -   picks a random entry of a list\n      can only be used as an answer\n" +
		"/play {game}\n   -   start playing a game\n"

	ff20HelpMsg = "Possible commands:\n" +
		"/m7\n   -   Some Action\n" +
		"/ff20\n   -   Another Action"

	lastFF20Post = ""
)

// HandleSettingsRequest handels messages concerning a users settings.
func HandleSettingsRequest(message whatsapp.TextMessage) {
	// Print Settings
	if strings.HasPrefix(strings.ToLower(message.Text), "/printsettings") { // TODO Update if needed

		settings := wabot.GetUserSettings(wabot.MessageToJid(message)).(Settings)

		infos := "@" + wabot.MessageToName(message) + "'s settings:"
		infos += "\nStdTranslationFrom: " + settings.StdTranslationFrom.String()
		infos += "\nStdTranslationTo: " + settings.StdTranslationTo.String()

		wabot.WriteTextMessage(infos, message.Info.RemoteJid)
	}
}

// HandlePickRequest picks one entry of a sent list and replies it
func HandlePickRequest(message whatsapp.TextMessage) {

	// Try to find the quoted message
	txt := ""
	if message.ContextInfo.QuotedMessage != nil {
		txt = message.ContextInfo.QuotedMessage.GetConversation()

		// Possible picks
		options := strings.Split(txt, "\n")

		// Reply a random entry
		rand.Seed(time.Now().UnixNano())
		wabot.WriteTextMessage("Random pick:\n"+options[rand.Intn(len(options))], message.Info.RemoteJid)
	} else {
		wabot.WriteTextMessage("Sorry!\nThis command can only be used while quoting a message", message.Info.RemoteJid)
	}

}
