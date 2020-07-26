package main

import (
	"strings"

	// "bufio"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	witai "github.com/wit-ai/wit-go"
	// "os"
	// witai "github.com/wit-ai/wit-go"
)

var machaaoAPIToken string = "3f79e9c0-c455-11ea-ad9a-094460ab21b1"
var witAiToken string = "HNNC6IHVUOVQUGH4ANQJJJHFQEQ326CG"

type memesResponse struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Nsfw      bool   `json:"nsfw"`
	Spoiler   bool   `json:"spoiler"`
}

func getJokeTagUsingWitAI(message string) string {
	client := witai.NewClient(witAiToken)
	// Use client.SetHTTPClient() to set custom http.Client

	msg, _ := client.Parse(&witai.MessageRequest{
		Query: message,
	})

	return msg.Entities["local_search_query"].([]interface{})[0].(map[string]interface{})["value"].(string)
}

func getJoke(tag string) string {

	var url string = "https://sv443.net/jokeapi/v2/joke/Any?format=txt&contains=" + tag

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Accept", "text/plain")
	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func simpleReply(userID []string, message string, apiToken string) {

	if strings.ToLower(message) == "random jokes" {
		message = getJoke("%20")
	} else if strings.ToLower(message) == "random memes" {
		title, url, subreddit := getMemes()

		_ = title
		body := map[string]interface{}{
			"identifier": "BROADCAST_FB_QUICK_REPLIES",
			"source":     "firebase",
			"users":      userID,
			"message": map[string]interface{}{
				"attachment": map[string]interface{}{
					"type": "template",
					"payload": map[string]interface{}{
						"template_type": "generic",
						"elements": []map[string]string{
							{
								"subtitle":  "Courtesy " + subreddit,
								"image_url": url,
							},
						},
					},
				},
				"quick_replies": []map[string]string{
					{
						"content_type": "text",
						"payload":      "Random Jokes",
						"title":        "Random Jokes",
					},
					{
						"content_type": "text",
						"payload":      "Random Memes",
						"title":        "Random Memes",
					},
				},
			},
		}
		log.Println("Sending Message to user")

		var urlMachaao string = "https://ganglia-dev.machaao.com/v1/messages/send"
		// var url string = "http://127.0.0.1:5000/upload"

		jsonValue, _ := json.Marshal(body)

		// fmt.Println(jsonValue)

		req, err := http.NewRequest("POST", urlMachaao, bytes.NewBuffer(jsonValue))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api_token", apiToken)

		fmt.Println(req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		bodyf, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(bodyf))

		return

	} else {
		var tag string = getJokeTagUsingWitAI(message)
		message = getJoke(tag)

		if message[:9] == "Error 106" {
			message = "Sorry, no jokes found"
		}
	}

	log.Println("Sending Message to user")

	var url string = "https://ganglia-dev.machaao.com/v1/messages/send"
	// var url string = "http://127.0.0.1:5000/upload"

	body := map[string]interface{}{
		"identifier": "BROADCAST_FB_QUICK_REPLIES",
		"source":     "firebase",
		"users":      userID,
		"message": map[string]interface{}{
			"text": message,
			"quick_replies": []map[string]string{
				{
					"content_type": "text",
					"payload":      "Random Jokes",
					"title":        "Random Jokes",
				},
				{
					"content_type": "text",
					"payload":      "Random Memes",
					"title":        "Random Memes",
				},
			},
		},
	}

	jsonValue, _ := json.Marshal(body)

	// fmt.Println(jsonValue)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_token", apiToken)

	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	bodyf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(bodyf))
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
		return []byte(machaaoAPIToken), nil
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

	if messageText == "hi" {
		quickReply(r.Header["User_id"], messageText, machaaoAPIToken)
	} else {
		simpleReply(r.Header["User_id"], messageText, machaaoAPIToken)
	}
}

func main() {
	port := getPort()
	http.HandleFunc("/machaao_hook", messageHandler)

	log.Println("[-] Listening on...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}

	// getMemes()
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4747"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}

func quickReply(userID []string, message string, apiToken string) {

	log.Println("Sending QR to user")

	var url string = "https://ganglia-dev.machaao.com/v1/messages/send"
	// var url string = "http://127.0.0.1:5000/upload"

	body := map[string]interface{}{
		"identifier": "BROADCAST_FB_QUICK_REPLIES",
		"source":     "firebase",
		"users":      userID,
		"message": map[string]interface{}{
			"text": "Hello, My name is Witty - Your funny friend ;)",
			"quick_replies": []map[string]string{
				{
					"content_type": "text",
					"payload":      "Random Jokes",
					"title":        "Random Jokes",
				},
				{
					"content_type": "text",
					"payload":      "Random Memes",
					"title":        "Random Memes",
				},
			},
		},
	}

	jsonValue, _ := json.Marshal(body)

	// fmt.Println(jsonValue)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_token", apiToken)

	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	bodyf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(bodyf))
}

func getMemes() (string, string, string) {

	var url string = "https://meme-api.herokuapp.com/gimme"

	var jsonBody memesResponse

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &jsonBody)

	return jsonBody.Title, jsonBody.URL, jsonBody.Subreddit
}
