package main

type Album struct {
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Rating     string `json:"rating"`
	ArtworkURL string `json:"artwork_url"`
}
