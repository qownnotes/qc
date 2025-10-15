package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/qownnotes/qc/config"
	"github.com/qownnotes/qc/entity"

	"net/url"

	"github.com/gorilla/websocket"
)

var (
	snippetCacheFile string
)

type Message struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type NoteFolderSwitchMessage struct {
	Token string `json:"token"`
	Type  string `json:"type"`
	Data  int    `json:"data"`
}

type ResultMessage struct {
	Type            string               `json:"type"`
	CommandSnippets []entity.SnippetInfo `json:"data"`
}

type NoteFolderResultMessage struct {
	Type        string                  `json:"type"`
	NoteFolders []entity.NoteFolderInfo `json:"data"`
	CurrentId   int                     `json:"currentId"`
}

type NoteFolderSwitchResultMessage struct {
	Type     string `json:"type"`
	Switched bool   `json:"data"`
}

// type CommandSnippet struct {
// 	Command      string `json:"command"`
// 	Description  string `json:"description"`
// 	Tags  		 []string `json:"tags"`
// }

func FetchSnippetsData() []entity.SnippetInfo {
	// log.Printf("Connecting to QOwnNotes on %s", u.String())

	initSnippetCacheFile()
	var snippetData []byte

	c, err := connectSocket()
	if err != nil {
		snippetData = readSnippetCacheFile()

		if snippetData == nil {
			log.Fatal(
				"Couldn't connect to QOwnNotes websocket, did you enable the socket server? Error: ",
				err,
			)
		} else {
			log.Println("Couldn't connect to QOwnNotes websocket, but found cached data in " + snippetCacheFile)
		}
	}

	if snippetData == nil {
		defer func() {
			if err := c.Close(); err != nil {
				log.Printf("Error closing connection: %v", err)
			}
		}()

		message := Message{
			Token: config.Conf.QOwnNotes.Token,
			Type:  "getCommandSnippets",
		}

		m, err := json.Marshal(message)
		if err != nil {
			log.Fatal("Couldn't marshal message: ", err)
		}

		if err = c.WriteMessage(websocket.TextMessage, m); err != nil {
			log.Fatal("Couldn't send command to QOwnNotes: ", err)
		}

		_, snippetData, err = c.ReadMessage()
		if err != nil {
			log.Fatalf("Couldn't read message from QOwnNotes: %v", err)
		}
	}

	var resultMessage ResultMessage
	err = json.Unmarshal(snippetData, &resultMessage)
	if err != nil {
		log.Fatalf("Couldn't interpret message from QOwnNotes: %v", err)
	}

	switch resultMessage.Type {
	case "tokenQuery":
		log.Fatal("Please execute \"qc configure\" and configure your token from QOwnNotes!")
	case "commandSnippets":
		writeSnippetCacheFile(snippetData)

		// log.Printf("CommandSnippets: %v", resultMessage.CommandSnippets)
		return resultMessage.CommandSnippets
	default:
		log.Fatal("Did not understand response from QOwnNotes!")
	}

	return []entity.SnippetInfo{}
}

func FetchNoteFolderData() (noteFolderInfo []entity.NoteFolderInfo, currentId int) {
	var noteFolderData []byte

	c, err := connectSocket()
	if err != nil {
		log.Fatal(
			"Couldn't connect to QOwnNotes websocket, did you enable the socket server? Error: ",
			err,
		)
	}

	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	message := Message{
		Token: config.Conf.QOwnNotes.Token,
		Type:  "getNoteFolders",
	}

	m, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Couldn't marshal message: ", err)
	}

	if err = c.WriteMessage(websocket.TextMessage, m); err != nil {
		log.Fatal("Couldn't send command to QOwnNotes: ", err)
	}

	_, noteFolderData, err = c.ReadMessage()
	if err != nil {
		log.Fatalf("Couldn't read message from QOwnNotes: %v", err)
	}

	var resultMessage NoteFolderResultMessage
	// log.Printf("Connecting to QOwnNotes on vs", noteFolderData)
	err = json.Unmarshal(noteFolderData, &resultMessage)
	if err != nil {
		log.Fatalf(
			"Couldn't interpret message from QOwnNotes: %v\nYou need at least QOwnNotes 22.7.1!",
			err,
		)
	}

	switch resultMessage.Type {
	case "tokenQuery":
		log.Fatal("Please execute \"qc configure\" and configure your token from QOwnNotes!")
	case "noteFolders":
		// log.Printf("NoteFolders: %v", resultMessage.NoteFolders)
		return resultMessage.NoteFolders, resultMessage.CurrentId
	default:
		log.Fatal("Did not understand response from QOwnNotes!")
	}

	return []entity.NoteFolderInfo{}, 0
}

func SwitchNoteFolder(id int) {
	c, err := connectSocket()
	if err != nil {
		log.Fatal(
			"Couldn't connect to QOwnNotes websocket, did you enable the socket server? Error: ",
			err,
		)
	}

	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	message := NoteFolderSwitchMessage{
		Token: config.Conf.QOwnNotes.Token,
		Type:  "switchNoteFolder",
		Data:  id,
	}

	m, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Couldn't marshal message: ", err)
	}

	if err = c.WriteMessage(websocket.TextMessage, m); err != nil {
		log.Fatal("Couldn't send command to QOwnNotes: ", err)
	}

	_, noteFolderData, err := c.ReadMessage()
	if err != nil {
		log.Fatalf("Couldn't read message from QOwnNotes: %v", err)
	}

	var resultMessage NoteFolderSwitchResultMessage
	err = json.Unmarshal(noteFolderData, &resultMessage)
	if err != nil {
		log.Fatalf(
			"Couldn't interpret message from QOwnNotes: %v\nYou need at least QOwnNotes 22.7.1!",
			err,
		)
	}

	switch resultMessage.Type {
	case "tokenQuery":
		log.Fatal("Please execute \"qc configure\" and configure your token from QOwnNotes!")
	case "switchedNoteFolder":
		if resultMessage.Switched {
			fmt.Println("Note folder was switched.")
		} else {
			fmt.Println("Note folder was not switched.")
		}
	default:
		log.Fatal("Did not understand response from QOwnNotes!")
	}
}

func connectSocket() (*websocket.Conn, error) {
	u := getSocketUrl()
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return c, err
}

func getSocketUrl() url.URL {
	u := url.URL{
		Scheme: "ws",
		Host:   "127.0.0.1:" + strconv.Itoa(config.Conf.QOwnNotes.WebSocketPort),
	}
	return u
}

func initSnippetCacheFile() {
	if snippetCacheFile == "" {
		dir, err := config.GetDefaultConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		snippetCacheFile = filepath.Join(dir, "snippets.cache")
	}
}

func writeSnippetCacheFile(data []byte) {
	if err := os.WriteFile(snippetCacheFile, data, 0600); err != nil {
		log.Fatal("Could not write snippet cache file: ", err)
	}
}

func readSnippetCacheFile() []byte {
	_, err := os.Stat(snippetCacheFile)

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	data, err := os.ReadFile(snippetCacheFile)

	if err != nil {
		log.Fatal("Could not read snippet cache file: ", err)
	}

	return data
}
