package main

import (
	"github.com/machaao/machaao-go"
)

// MemesResponse a
type MemesResponse struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Nsfw      bool   `json:"nsfw"`
	Spoiler   bool   `json:"spoiler"`
}

// UserTags s
type UserTags struct {
	ID          string `json:"_id"`
	DisplayName string `json:"displayName"`
	Values      []bool `json:"values"`
	Name        string `json:"name"`
}

// MemeSubreddits Meme keywords mapped to subreddits.
var MemeSubreddits = map[string]string{
	"school":      "gradschoolmemes",
	"college":     "gradschoolmemes",
	"programming": "ProgrammerHumor",
	"photoshop":   "photoshopbattles",
	"politics":    "PresidentialRaceMemes",
	"nsfw":        "NSFWMeme",
}

// NsfwSubreddits h
var NsfwSubreddits = []string{"NSFWFunny", "NSFWMeme", "MemesNSFW", "Nsfwhumour", "nsfw", "NSFW_GIF"}

// RndSubreddits a
var RndSubreddits = []string{"memes", "dankmemes", "Memes_Of_The_Dank", "ComedyCemetery", "FellowKids", "wholesomememes", "ProtectAndServe"}

func main() {

	machaao.Server(MessageHandler)
}
