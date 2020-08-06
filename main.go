package main

import (
	"strings"

	// "bufio"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/machaao/machaao-go"
	witai "github.com/wit-ai/wit-go"
	// "os"
	// witai "github.com/wit-ai/wit-go"
)

type memesResponse struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Nsfw      bool   `json:"nsfw"`
	Spoiler   bool   `json:"spoiler"`
}

var memeSubreddits = map[string]string{
	"school":     "gradschoolmemes",
	"college":    "gradschoolmemes",
	"university": "gradschoolmemes",
	"photoshop":  "photoshopbattles",
	"no context": "nocontextpics",
	"animals":    "AdviceAnimals",
	"nsfw":       "NSFWMeme",
}

var isUserAdult bool = false

func main() {
	machaao.Server(messageHandler)
}

func getMemeTagUsingWitAI(message string) string {
	client := witai.NewClient(machaao.WitAPIToken)
	// Use client.SetHTTPClient() to set custom http.Client

	msg, _ := client.Parse(&witai.MessageRequest{
		Query: message,
	})

	return msg.Entities["local_search_query"].([]interface{})[0].(map[string]interface{})["value"].(string)
}

func sendRandMeme(userID []string, message string) {

	log.Println("Sending Message to user")

	_, url, postlink := getMemes("")

	resp, err := machaao.SendMessage(getMemeBody(userID, url, postlink))

	if err != nil {
		log.Println(err)
	}

	log.Printf("SR POST Request Response %s", resp.Status)

}

func messageHandler(w http.ResponseWriter, r *http.Request) {

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

	fmt.Println(messageData)
	fmt.Println(messageText)

	fmt.Println(r.Header["User_id"])

	if strings.ToLower(messageText) == "hi" {
		quickReply(r.Header["User_id"], messageText)
	} else if strings.ToLower(messageText) == "random memes" || strings.ToLower(messageText) == "random meme" {
		sendRandMeme(r.Header["User_id"], messageText)
	} else if strings.ToLower(messageText) == "nsfw" {
		if isUserAdult {
			sendSpecificMemes(r.Header["User_id"], messageText)
		} else {
			checkAdultPrompt(r.Header["User_id"])
		}
	} else if messageText == "setADULT18" {
		setAdultVar(r.Header["User_id"])
	} else {
		sendSpecificMemes(r.Header["User_id"], messageText)
	}
}

func checkAdultPrompt(userID []string) {

	body := map[string]interface{}{
		"users": userID,
		"message": map[string]interface{}{
			"text": "Are you over 18 year old?",
			"quick_replies": []map[string]string{
				{
					"content_type": "text",
					"payload":      "setADULT18",
					"title":        "Yes, I'm over 18",
				},
				{
					"content_type": "text",
					"payload":      "no",
					"title":        "No",
				},
			},
		},
	}

	resp, err := machaao.SendMessage(body)

	if err != nil {
		log.Println(err)
	}

	log.Printf("Check Adult Prompt %s", resp.Status)

}

func setAdultVar(userID []string) {
	isUserAdult = true
	sendSpecificMemes(userID, "nsfw")
	log.Printf("NOW %s is set to ADULT", userID)
}

func quickReply(userID []string, message string) {

	log.Println("Sending QR to user")

	body := map[string]interface{}{
		"users": userID,
		"message": map[string]interface{}{
			"text": "Hello, My name is Witty - Your meme friend ;)",
			"quick_replies": []map[string]string{
				{
					"content_type": "text",
					"payload":      "Random Memes",
					"title":        "🙃 Random Memes",
				},
				{
					"content_type": "text",
					"payload":      "school",
					"title":        "School",
				},
				{
					"content_type": "text",
					"payload":      "photoshop",
					"title":        "Photoshop",
				},
				{
					"content_type": "text",
					"payload":      "no context",
					"title":        "No Context",
				},
				{
					"content_type": "text",
					"payload":      "nsfw",
					"title":        "NSFW",
				},
			},
		},
	}

	resp, err := machaao.SendMessage(body)

	if err != nil {
		log.Println(err)
	}

	log.Printf("QR POST Request Response %s", resp.Status)

}

func sendSpecificMemes(userID []string, message string) {

	var url, postlink string = "", ""
	if subreddit, ok := memeSubreddits[message]; ok {
		_, url, postlink = getMemes(subreddit)
	} else {
		_, url, postlink = getMemes("")
	}

	resp, err := machaao.SendMessage(getMemeBody(userID, url, postlink))

	if err != nil {
		log.Println(err)
	}

	log.Printf("Specific Meme POST Request Response %s", resp.Status)

}

func getMemes(subreddit string) (string, string, string) {

	var url string = ""

	if subreddit == "" {
		url = "https://meme-api.herokuapp.com/gimme"
	} else {
		url = "https://meme-api.herokuapp.com/gimme/" + subreddit
	}

	var jsonBody memesResponse

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &jsonBody)

	return jsonBody.Title, jsonBody.URL, jsonBody.PostLink
}

func getMemeBody(userID []string, url string, postlink string) interface{} {
	body := map[string]interface{}{
		"users": userID,
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "generic",
					"elements": []map[string]interface{}{
						{
							"image_url": url,
							"buttons": []map[string]string{
								{
									"type":  "web_url",
									"url":   postlink,
									"title": "ℹ️ Source",
								},
							},
						},
					},
				},
			},
			"quick_replies": []map[string]string{
				{
					"content_type": "text",
					"payload":      "Random Memes",
					"title":        "🙃 Random Memes",
				},
				{
					"content_type": "text",
					"payload":      "school",
					"title":        "School",
				},
				{
					"content_type": "text",
					"payload":      "photoshop",
					"title":        "Photoshop",
				},
				{
					"content_type": "text",
					"payload":      "no context",
					"title":        "No Context",
				},
				{
					"content_type": "text",
					"payload":      "nsfw",
					"title":        "NSFW",
				},
			},
		},
	}

	return body
}
