package main

import (
	// "bufio"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	// "os"
	// witai "github.com/wit-ai/wit-go"
)

func getJoke() string {
	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)

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

	if message == "joke" {
		message = getJoke()
	}

	log.Println("Sending Message to user")

	var url string = "https://ganglia-dev.machaao.com/v1/messages/send"
	// var url string = "http://127.0.0.1:5000/upload"

	body := map[string]interface{}{"identifier": "BROADCAST_FB_QUICK_REPLIES", "source": "firebase", "users": userID, "message": map[string]string{"text": message}}

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
	fmt.Println("response Headers:", resp.Header)
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
		return []byte("3f79e9c0-c455-11ea-ad9a-094460ab21b1"), nil
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

	simpleReply(r.Header["User_id"], messageText, "3f79e9c0-c455-11ea-ad9a-094460ab21b1")
}

func main() {
	// client := witai.NewClient("HNNC6IHVUOVQUGH4ANQJJJHFQEQ326CG")
	// // Use client.SetHTTPClient() to set custom http.Client

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Enter message to send.")
	// query, _ := reader.ReadString('\n')

	// msg, _ := client.Parse(&witai.MessageRequest{
	// 	Query: query,
	// })
	// fmt.Printf("%v\n", msg)
	http.HandleFunc("/machaao_hook", messageHandler)

	fmt.Printf("Starting server at http://127.0.0.1:8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
