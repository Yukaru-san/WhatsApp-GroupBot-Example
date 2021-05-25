package main

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	wabot "github.com/Yukaru-san/WhatsApp-GroupBot"
)

// TODO: VERSUS MODE, CHANGE WORD-SET

// TODO: Use this to make hangman work for each group once
type hangmanInstance struct {
	players          []string // Registered players (Jid's)
	hangmanIsRunning bool
	hangmanIsPregame bool

	endgameTimerInterval int64
	endgameTimer         int64 // game ends when the timer here is reached

	currentHangmanState int
	currentWord         string
	currentlyKnownWord  string
	alreadyTriedChars   []string
}

var (
	players          []string // Registered players (Jid's)
	hangmanIsRunning = false
	hangmanIsPregame = false

	endgameTimerInterval = (int64(time.Minute.Seconds()) * 10)
	endgameTimer         int64 // game ends when the timer here is reached

	currentHangmanState = -1
	currentWord         = ""
	currentlyKnownWord  = ""
	alreadyTriedChars   []string

	hangmanStates = []string{
		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"``` _|_        ```\n" +
			"```|   |______```\n" +
			"```|          |```\n" +
			"```|__________|```\n",

		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |    o    ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"``` _|_        ```\n" +
			"```|   |______ ```\n" +
			"```|          |```\n" +
			"```|__________|```",

		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |    o    ```\n" +
			"```  |   /     ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"``` _|_        ```\n" +
			"```|   |______ ```\n" +
			"```|          |```\n" +
			"```|__________|```",

		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |    o    ```\n" +
			"```  |   /|    ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"``` _|_        ```\n" +
			"```|   |______ ```\n" +
			"```|          |```\n" +
			"```|__________|```",

		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |    o    ```\n" +
			"```  |   /|\\   ```\n" +
			"```  |         ```\n" +
			"```  |         ```\n" +
			"``` _|_        ```\n" +
			"```|   |______ ```\n" +
			"```|          |```\n" +
			"```|__________|```",

		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |    o    ```\n" +
			"```  |   /|\\   ```\n" +
			"```  |    |    ```\n" +
			"```  |         ```\n" +
			"``` _|_        ```\n" +
			"```|   |______ ```\n" +
			"```|          |```\n" +
			"```|__________|```",

		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |    o    ```\n" +
			"```  |   /|\\   ```\n" +
			"```  |    |    ```\n" +
			"```  |   /     ```\n" +
			"``` _|_        ```\n" +
			"```|   |______ ```\n" +
			"```|          |```\n" +
			"```|__________|```",

		"```   ____     ```\n" +
			"```  |    |    ```\n" +
			"```  |    o    ```\n" +
			"```  |   /|\\   ```\n" +
			"```  |    |    ```\n" +
			"```  |   / \\   ```\n" +
			"``` _|_        ```\n" +
			"```|   |______ ```\n" +
			"```|          |```\n" +
			"```|__________|```",
	}
)

// Sets all settings to default
func resetSettings(keepPlayers bool) {
	updateGameTimer()
	if !keepPlayers {
		players = nil
		players = []string{}
	}
	alreadyTriedChars = nil
	alreadyTriedChars = []string{}
	hangmanIsRunning = false
	hangmanIsPregame = true
	currentHangmanState = -1
	currentWord = ""
	currentlyKnownWord = ""
}

// InitiateHangman initiates the game
func InitiateHangman(message whatsapp.TextMessage) {
	// Setup
	wabot.WriteTextMessage("OK!\n\nTo register for the game, write \"/join\"\n\nTo start the game, write \"/start\"\n\nTo end the game, write at any point \"/stopgame\"", message.Info.RemoteJid)

	updateGameTimer()
	resetSettings(false)

	// Handle messages
	wabot.SetDefaultTextHandleFunction(func(message whatsapp.TextMessage) {
		HandleHangmanMessage(message)
	})

	// End the game after a given amount of time
	go (func() {
		for hangmanIsRunning {
			if endgameTimer < time.Now().Unix() {
				wabot.WriteTextMessage("Hangman:\nSince nobody played in the last 10 minutes the game got canceled.", message.Info.RemoteJid)
				wabot.SetDefaultTextHandleFunction(func(whatsapp.TextMessage) {})
				resetSettings(false)
			}
		}
	})()
}

// HandleHangmanMessage is for everything needed to play
func HandleHangmanMessage(message whatsapp.TextMessage) {
	// Pregame commands
	if hangmanIsPregame {

		if strings.HasPrefix(strings.ToLower(message.Text), "/join") {
			isRegistered := isPlayerRegistered(wabot.MessageToJid(message))
			if isRegistered {
				wabot.WriteTextMessage("Hangman:\nYou are already registered @"+wabot.MessageToName(message)+"!", message.Info.RemoteJid)
			} else {
				players = append(players, wabot.MessageToJid(message))
				wabot.WriteTextMessage("Hangman:\nSuccessfully registered @"+wabot.MessageToName(message)+"!", message.Info.RemoteJid)
			}

		} else if strings.HasPrefix(strings.ToLower(message.Text), "/start") {
			startActualGame(message)

		} else if strings.HasPrefix(strings.ToLower(message.Text), "/quit") {
			removePlayer(message)

			if len(players) == 0 {
				stopGameEmptyLobby()
				wabot.WriteTextMessage("Hangman:\nThe game was closed since there are no players left to play!", message.Info.RemoteJid)
			}

		} else if strings.HasPrefix(strings.ToLower(message.Text), "/stopgame") {
			if stopGame(message) {
				wabot.WriteTextMessage("Hangman:\nGame finished!", message.Info.RemoteJid)
			}
		}

		updateGameTimer()
		return
	}

	// Actual game commands
	if strings.ToLower(message.Text) == "/stopgame" {
		if stopGame(message) {
			wabot.WriteTextMessage("Hangman:\nGame got ended early!", message.Info.RemoteJid)
		}
	} else if isPlayerRegistered(wabot.MessageToJid(message)) {
		handlePlayerMessage(message)
	}

	updateGameTimer()
	return
}

// Player interaction with the game
func handlePlayerMessage(message whatsapp.TextMessage) {
	// User wrote the full & correct word
	if strings.ToLower(message.Text) == strings.ToLower(currentWord) {
		wabot.WriteTextMessage("Hangman:\n@"+wabot.MessageToName(message)+" has guessed the word!", message.Info.RemoteJid)
		roundFinished(message.Info.RemoteJid)

	} else if len(message.Text) == 1 { // Ignore things people might write to themselves
		if charWasAlreadyTried(strings.ToLower(message.Text)) {
			wabot.WriteTextMessage("Hangman:\nThis letter was tested already!", message.Info.RemoteJid)

		} else if strings.Contains(strings.ToLower(currentWord), strings.ToLower(message.Text)) {
			// Check if the char is included
			correctInput(message)

		} else {
			// Wrong input
			wrongInput(message)
		}
	}
}

// Handles correct inputs
func correctInput(message whatsapp.TextMessage) {

	newlyKnownWord := ""
	newWordLen := 0

	// Adding the letters to the word
	for i := 0; i < len(currentWord); i++ {
		if i < len(currentWord)-1 {
			if strings.ToLower(string(currentWord[i])) == strings.ToLower(message.Text) || charWasAlreadyTried(strings.ToLower(string(currentWord[i]))) {
				newlyKnownWord += (string(currentWord[i]) + " ")
				newWordLen++
			} else {
				newlyKnownWord += ("_ ")
			}
		} else {
			if strings.ToLower(string(currentWord[i])) == strings.ToLower(message.Text) || charWasAlreadyTried(strings.ToLower(string(currentWord[i]))) {
				newlyKnownWord += (string(currentWord[i]))
				newWordLen++
			} else {
				newlyKnownWord += ("_")
			}
		}
	}

	// Won the Game
	if newWordLen == len(currentWord) {

		wabot.WriteTextMessage("Hangman:\n@"+wabot.MessageToName(message)+" has guessed the word *"+currentWord+"*!", message.Info.RemoteJid)
		roundFinished(message.Info.RemoteJid)

		// Game is not won yet
	} else {
		currentlyKnownWord = newlyKnownWord

		alreadyTriedChars = append(alreadyTriedChars, strings.ToLower(message.Text))
		wabot.WriteTextMessage("Hangman:\nThis letter is contained!\n\nYou already know:\n"+currentlyKnownWord+"\n\nTried letters:\n"+getReadableListOfTriedChars(), message.Info.RemoteJid)
	}
}

// Handles incorrect inputs
func wrongInput(message whatsapp.TextMessage) {

	// Too many errors, game end
	if currentHangmanState == 6 {

		wabot.WriteTextMessage(hangmanStates[currentHangmanState+1]+"\n\nYou lost the game!\nThe word was:\n*"+currentWord+"*", message.Info.RemoteJid)
		roundFinished(message.Info.RemoteJid)

	} else {
		// Game is still in progress, but there is one additional error

		currentHangmanState++

		alreadyTriedChars = append(alreadyTriedChars, strings.ToLower(message.Text))
		wabot.WriteTextMessage(hangmanStates[currentHangmanState]+"\n\nThis letter was *not* included! \n\nYou already know:\n"+currentlyKnownWord+"\n\nTried letters:\n"+getReadableListOfTriedChars(), message.Info.RemoteJid)
	}
}

// Returns a single string containing every tried char
func getReadableListOfTriedChars() string {
	list := ""
	for i := 0; i < len(alreadyTriedChars); i++ {
		list += (alreadyTriedChars[i] + " | ")
	}
	return list
}

func charWasAlreadyTried(char string) bool {
	alreadyTried := false
	for i := 0; i < len(alreadyTriedChars); i++ {
		if alreadyTriedChars[i] == char {
			alreadyTried = true
		}
	}
	return alreadyTried
}

// Actually sets the game into it's "ingame" mode
func startActualGame(msg whatsapp.TextMessage) {
	if isPlayerRegistered(wabot.MessageToJid(msg)) {

		// Read the words-file TODO add more possible Word-Lists and change this line:
		fileBytes, _ := ioutil.ReadFile("./hangmanWordlists/basic.txt")
		fileString := strings.ReplaceAll(string(fileBytes), "\r", "")

		words := strings.Split(fileString, "\n")

		// Pick one randomly
		rand.Seed(time.Now().UnixNano())
		currentWord = words[rand.Intn(len(words))]

		println("New Word:\n" + currentWord)

		// Disguise the word
		currentlyKnownWord = ""
		for i := 0; i < len(currentWord); i++ {
			// Easier to read
			if i < len(currentWord)-1 {
				currentlyKnownWord += "_ "
			} else {
				currentlyKnownWord += "_"
			}
		}

		// Set states
		hangmanIsPregame = false
		hangmanIsRunning = true
		wabot.WriteTextMessage("Hangman:\nGame started! Every message of participating players count as a try!\n\nThe word is...\n"+currentlyKnownWord, msg.Info.RemoteJid)
	} else {
		wabot.WriteTextMessage("Hangman:\nOnly registered users can start the game!", msg.Info.RemoteJid)
	}
}

// Transition into post-game
func roundFinished(jid string) {
	resetSettings(true)

	wabot.WriteTextMessage("To register for the next game, please type \"/join\"\n\nTo undo the registration, type \"/quit\"\n\nTo start the game, please write \"/start\"\n\nTo stop the game at any point, type \"/stopgame\"", jid)
}

// Stops the entire game
func stopGame(message whatsapp.TextMessage) bool {
	if isPlayerRegistered(wabot.MessageToJid(message)) {
		resetSettings(false)

		// Stop handling messages
		wabot.SetDefaultTextHandleFunction(func(message whatsapp.TextMessage) {
		})

		return true
	}

	wabot.WriteTextMessage("Hangman:\nOnly registered users can end the game!", message.Info.RemoteJid)
	return false
}

// Stops the entire game
func stopGameEmptyLobby() {
	resetSettings(false)

	// Stop handling messages
	wabot.SetDefaultTextHandleFunction(func(message whatsapp.TextMessage) {
	})
}

// Check if player is already registered
func isPlayerRegistered(jid string) bool {
	isRegistered := false
	for i := 0; i < len(players); i++ {
		if players[i] == jid {
			isRegistered = true
		}
	}
	return isRegistered
}

// Remove player from the list
func removePlayer(message whatsapp.TextMessage) {

	jid := wabot.MessageToJid(message)

	var temp []string
	for i := 0; i < len(players); i++ {
		if players[i] != jid {
			temp = append(temp, players[i])
		}
	}
	players = temp

	wabot.WriteTextMessage("Hangman:\n"+wabot.JidToName(jid)+" has left the game!", message.Info.RemoteJid)
}

// Updates the timer until the game automatically closes
func updateGameTimer() {
	endgameTimer = time.Now().Unix() + endgameTimerInterval
}
