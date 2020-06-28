package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/lithdew/youtube"
	"github.com/zmb3/spotify"
)

func getTNDAlbums(spo spotify.Client) []Album {
	res, err := youtube.LoadPlaylist(tndLoveList, 0)
	if err != nil {
		log.Printf("%s", err)
	}

	var aa []Album
	for _, vid := range res.Items {
		var alb Album
		if titleArtist := strings.Split(strings.Replace(vid.Title, " ALBUM REVIEW", "", -1), " - "); len(titleArtist) == 2 {
			alb.Title = titleArtist[0]
			alb.Artist = titleArtist[1]
		}

		desc := vid.Description

		re, _ := regexp.Compile(`\n(\d+\/10)`)
		if matches := re.FindStringSubmatch(desc); len(matches) > 0 {
			alb.Rating = matches[1]
		}

		aa = append(aa, alb)
	}

	updateAlbumsArtwork(spo, aa)

	return aa
}

func createRatingAlbums(aa []Album) map[int][]Album {
	ratingAlbums := make(map[int][]Album)
	for _, alb := range aa {
		if spl := strings.Split(alb.Rating, "/"); len(spl) == 2 {
			if n, err := strconv.Atoi(spl[0]); err == nil {
				ratingAlbums[n] = append(ratingAlbums[n], alb)
			}
		}
	}

	return ratingAlbums
}

func updateAlbumsArtwork(spo spotify.Client, aa []Album) {
	for i := 0; i < len(aa); i++ {
		var artworkURL string
		exponentialBackoff(func() error {
			var err error
			artworkURL, err = getAlbumArtworkURL(spo, aa[i])
			return err
		}, 4)
		aa[i].ArtworkURL = artworkURL
	}
}

func getAlbumArtworkURL(spo spotify.Client, a Album) (string, error) {
	res, err := spo.Search(a.Title+" "+a.Artist, spotify.SearchTypeAlbum)
	if err != nil {
		return "", err
	}

	var imgURL string
	if aa := res.Albums.Albums; len(aa) > 0 {
		if imgs := aa[0].Images; len(imgs) > 0 {
			imgURL = imgs[0].URL
		}
	}

	return imgURL, nil
}
