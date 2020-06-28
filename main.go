package main

import (
	"context"
	"log"
	"os"

	"time"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify"
)

const tndLoveList = "PLP4CSgl7K7oo93I49tQa0TLB8qY3u7xuO"
const spotifyClientID = "d28c4c48ea64411bb48f989aa1f60c02"
const subMUCoreSpotifyPlaylistID = "1SCKNtsjWFtQJC6bLu4Kv9"

var spotifySecretKey string

var ratingAlbums = make(map[int][]Album)
var subMuCoreAlbums = []Album{}

func init() {
	spotifySecretKey = os.Getenv("SPOTIFY_SECRET_KEY")
}

func main() {
	spo := newSpotifyClient()
	subMuCoreAlbums = getSubMUCoreAlbums(spo)
	ratingAlbums = createRatingAlbums(getTNDAlbums(spo))

	e := newEcho()
	e.Logger.Fatal(e.Start(":8080"))
}

func newSpotifyClient() spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     spotifyClientID,
		ClientSecret: spotifySecretKey,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	return spotify.Authenticator{}.NewClient(token)
}

func exponentialBackoff(fn func() error, retries int) {
	expBackoff := 0
	for i := 0; i < retries; i++ {
		time.Sleep(time.Second * time.Duration(expBackoff))
		err := fn()
		if err == nil {
			break
		}

		if expBackoff == 0 {
			expBackoff = 1
		} else {
			expBackoff *= 2
		}
		log.Printf("failed executing function: %s. trying again in %v seconds", err, expBackoff)
	}
}
