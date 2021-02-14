package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/machaao/machaao-go"
)

// MessageHandler Handles upcoming and outgoing messages.
func MessageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var bodyData string = string(body)
	var tokenString string = bodyData[8:(len(bodyData) - 2)]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(machaao.MachaaoAPIToken), nil
	})

	_ = token

	if err != nil {
		fmt.Println(err)
	}

	messageData := claims["sub"].(map[string]interface{})["messaging"].([]interface{})[0].(map[string]interface{})["message_data"]
	messageText := messageData.(map[string]interface{})["text"].(string)

	actionType := messageData.(map[string]interface{})["action_type"]
	timezone := claims["sub"].(map[string]interface{})["messaging"].([]interface{})[0].(map[string]interface{})["user"].(map[string]interface{})["timezone"]

	var userID []string = r.Header["User_id"]

	log.Printf("User: %s, Timezone: %s, ActionType: %s, Message: %s", userID[0], timezone, actionType, messageText)

	// Send quick reply
	if strings.ToLower(messageText) == "hi" {

		QuickReply(userID)

		// Check if user ask for random memes
	} else if strings.ToLower(messageText) == "random memes" || strings.ToLower(messageText) == "random meme" {

		SendRandMeme(userID, messageText)

		// Check if user ask for nsfw content
	} else if strings.ToLower(messageText) == "nsfw" {

		// Check if user is 18+, send meme, else check in DATABASE or send ADULT CHECK prompt to new user.
		resp, _ := machaao.GetUserTag(userID[0])

		var tagData []interface{}

		body1, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body1, &tagData)

		if len(tagData) == 0 {

			CheckAdultPrompt(userID)

		} else if tagData[0].(map[string]interface{})["name"] == "adult" {

			log.Printf("%s is an ADULT", userID)

			// Get random meme subreddit.
			rndNum1 := rand.Intn(len(NsfwSubreddits))
			SendNSFWMemes(userID, NsfwSubreddits[rndNum1])

		} else {

			CheckAdultPrompt(userID)

		}

	} else if messageText == "setADULT18" {
		SetAdultVar(userID)
	} else {
		SendSpecificMemes(userID, messageText, messageText)
	}
}
