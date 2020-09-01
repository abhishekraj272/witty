package main

import (
	"log"
	"math/rand"

	"github.com/machaao/machaao-go"
)

// SendRandMeme Sends random meme to user.
func SendRandMeme(userID []string, message string) {

	log.Println("Sending Message to user")

	rndNum := rand.Intn(len(RndSubreddits))

	_, url, postlink := GetMemes(RndSubreddits[rndNum], userID)

	resp, err := machaao.SendMessage(GetMemeBody(userID, url, postlink, "Random Memes"))

	if err != nil {
		log.Println(err)
	}

	log.Printf("SR POST Request Response %s", resp.Status)

}

// SendNSFWMemes Sends NSFW memes to user.
func SendNSFWMemes(userID []string, subreddit string) {

	_, url, postlink := GetMemes(subreddit, userID)

	resp, err := machaao.SendMessage(GetMemeBody(userID, url, postlink, "nsfw"))

	if err != nil {
		log.Println(err)
	}

	log.Printf("Specific Meme POST Request Response %s", resp.Status)

}

// QuickReply Sends quickreply to user.
func QuickReply(userID []string) {

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

	resp, err := machaao.SendMessage(body)

	if err != nil {
		log.Println(err)
	}

	log.Printf("QR POST Request Response %s", resp.Status)

}

// SendSpecificMemes Sends some specific memes to user.
func SendSpecificMemes(userID []string, message string, memeType string) {

	var url, postlink string = "", ""
	if subreddit, ok := MemeSubreddits[message]; ok {
		_, url, postlink = GetMemes(subreddit, userID)
	} else {
		_, url, postlink = GetMemes("", userID)
	}

	resp, err := machaao.SendMessage(GetMemeBody(userID, url, postlink, memeType))

	if err != nil {
		log.Println(err)
	}

	log.Printf("Specific Meme POST Request Response %s", resp.Status)

}
