package main

import (
	"fmt"
	"strings"

	"github.com/Rhymen/go-whatsapp"
	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
	"github.com/bregydoc/gtranslate"
	"golang.org/x/text/language"
)

// HandleTranslateRequest handles a translation request
func HandleTranslateRequest(message whatsapp.TextMessage) {

	settings := wabot.GetUserSettings(wabot.MessageToJid(message)).(Settings)

	langFrom := settings.StdTranslationFrom
	langTo := settings.StdTranslationTo

	commandPretext := 2

	// If it was /tf or /tt
	if CheckForLanguageChange(message) {
		return
	}

	// If a translation direction was given use it instead
	if strings.HasPrefix(message.Text, "/t-") || strings.HasPrefix(message.Text, "!t-") {
		// Splitting
		wantedLangFromString := message.Text[3:5]
		wantedLangToString := message.Text[6:8]

		// Try to parse it
		wantedLangFrom, err := language.Parse(wantedLangFromString)
		wantedLangTo, err := language.Parse(wantedLangToString)

		// Success!
		if err == nil {
			langFrom = wantedLangFrom
			langTo = wantedLangTo
			commandPretext = 9 // Start at the end of the command
		} else {
			fmt.Println("Error during translation. ", err.Error())
		}
	}

	var txt string

	// If the message quotes another one, take that one
	if message.ContextInfo.QuotedMessage != nil {
		txt = message.ContextInfo.QuotedMessage.GetConversation()
	} else if len(message.Text) > 2 { // Else, use the given one
		txt = message.Text[commandPretext:] // Ignore the command-part
	}

	// Translate the message
	if len(txt) > 0 {
		tl, _ := gtranslate.Translate(txt, langFrom, langTo)
		wabot.WriteTextMessage(tl, message.Info.RemoteJid)
	}
}

// CheckForLanguageChange handles /tf and /tt
func CheckForLanguageChange(message whatsapp.TextMessage) bool {
	// Translate From xx
	if strings.HasPrefix(message.Text, "/tf") {
		newTranslation, err := language.Parse(message.Text[4:])

		// Input error
		if err != nil {
			wabot.WriteTextMessage("Sorry!\n Couldn't parse your Input.\nExample: /tf de -> translate from german", message.Info.RemoteJid)
			return true
		}

		// Search User and change the entry
		wabot.AddUserByJid(wabot.MessageToJid(message))

		// Change settings
		settings := wabot.GetUserSettings(wabot.MessageToJid(message)).(Settings)
		settings.StdTranslationFrom = newTranslation
		wabot.ChangeUserSettings(wabot.MessageToJid(message), settings)

		// Send a reply
		wabot.WriteTextMessage("Changed @"+wabot.MessageToName(message)+"'s {From} to "+newTranslation.String(), message.Info.RemoteJid)
		return true
	}

	// Translate to xx
	if strings.HasPrefix(message.Text, "/tt") {
		newTranslation, err := language.Parse(message.Text[4:])

		// Input error
		if err != nil {
			wabot.WriteTextMessage("Sorry!\n Couldn't parse your Input.\nExample: /tt eng -> translate to english", message.Info.RemoteJid)
			return true
		}

		// Search User and change the entry
		wabot.AddUserByJid(wabot.MessageToJid(message))

		// Change settings
		settings := wabot.GetUserSettings(wabot.MessageToJid(message)).(Settings)
		settings.StdTranslationTo = newTranslation
		wabot.ChangeUserSettings(wabot.MessageToJid(message), settings)

		// Send a reply
		wabot.WriteTextMessage("Changed @"+wabot.MessageToName(message)+"'s {To} to "+newTranslation.String(), message.Info.RemoteJid)
		return true
	}

	// Translate From AND To
	if strings.HasPrefix(message.Text, "/ts-") {
		// Splitting
		newTranslationFrom, err := language.Parse(message.Text[4:6])
		newTranslationTo, err := language.Parse(message.Text[7:9])

		// Input error
		if err != nil {
			wabot.WriteTextMessage("Sorry!\n Couldn't parse your Input.\nExample: /tt eng -> translate to english", message.Info.RemoteJid)
			return true
		}

		// Change settings FROM
		settingsFrom := wabot.GetUserSettings(wabot.MessageToJid(message)).(Settings)
		settingsFrom.StdTranslationFrom = newTranslationFrom
		wabot.ChangeUserSettings(wabot.MessageToJid(message), settingsFrom)

		// Change settings TO
		settingsTo := wabot.GetUserSettings(wabot.MessageToJid(message)).(Settings)
		settingsTo.StdTranslationTo = newTranslationTo
		wabot.ChangeUserSettings(wabot.MessageToJid(message), settingsTo)

		// Send a reply
		wabot.WriteTextMessage("Changed @"+wabot.MessageToName(message)+"'s\n{From} to "+newTranslationFrom.String()+"\n{To} to "+newTranslationTo.String(), message.Info.RemoteJid)
		return true
	}

	return false
}
