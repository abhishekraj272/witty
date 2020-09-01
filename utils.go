package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/machaao/machaao-go"
	witai "github.com/wit-ai/wit-go"
)

// GetMemes Fetches meme from meme api https://github.com/R3l3ntl3ss/Meme_Api
func GetMemes(subreddit string, userID []string) (string, string, string) {

	var url string = ""

	if subreddit == "" {
		url = "https://meme-api.herokuapp.com/gimme"
	} else {
		url = "https://meme-api.herokuapp.com/gimme/" + subreddit
	}

	var jsonBody MemesResponse

	req, err1 := http.NewRequest("GET", url, nil)

	if err1 != nil {
		log.Println(err1)
		QuickReply(userID)
	}

	client := &http.Client{}
	resp, err2 := client.Do(req)

	if err2 != nil {
		log.Println(err2)
		QuickReply(userID)
	}

	body, err3 := ioutil.ReadAll(resp.Body)

	if err3 != nil {
		log.Println(err3)
		QuickReply(userID)
	}

	json.Unmarshal(body, &jsonBody)

	return jsonBody.Title, jsonBody.URL, jsonBody.PostLink
}

// CheckAdultPrompt Asks user if s/he is an ADULT.
func CheckAdultPrompt(userID []string) {

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

// GetMemeBody Creates a interface to be sent to send message API.
func GetMemeBody(userID []string, url string, postlink string, memeType string) interface{} {
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
									"type":    "postback",
									"payload": memeType,
									"title":   "🔁 Repeat",
								},
								{
									"type":  "web_url",
									"url":   postlink,
									"title": "👀 View",
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
					"payload":      "politics",
					"title":        "Politics",
				},
				{
					"content_type": "text",
					"payload":      "programming",
					"title":        "Programming",
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

// SetAdultVar Tag user with ADULT tag.
func SetAdultVar(userID []string) {

	body := map[string]interface{}{
		"tag":         "adult",
		"source":      "web",
		"status":      1,
		"displayName": "Adult",
	}

	machaao.TagUser(userID[0], body)

	SendSpecificMemes(userID, "nsfw", "nsfw")
	log.Printf("NOW %s is set to ADULT", userID)
}

func getMemeTagUsingWitAI(message string) string {
	client := witai.NewClient(machaao.WitAPIToken)
	// Use client.SetHTTPClient() to set custom http.Client

	msg, _ := client.Parse(&witai.MessageRequest{
		Query: message,
	})

	return msg.Entities["local_search_query"].([]interface{})[0].(map[string]interface{})["value"].(string)
}
