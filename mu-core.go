package main

import (
	"log"
	"strings"

	"github.com/zmb3/spotify"
	"go.bmvs.io/orderedset"
)

func getSubMUCoreAlbums(spo spotify.Client) []Album {
	lim := 100

	res, err := spo.GetPlaylistTracks(subMUCoreSpotifyPlaylistID)
	if err != nil {
		log.Printf("error fetching sub-/mu/core tracks from spotify = %s\n", err)
		return []Album{}
	}

	totalTracks := res.Total
	tt := []spotify.PlaylistTrack{}
	for offset := 0; offset < totalTracks; offset += 100 {
		fn := func() error {
			res, err = spo.GetPlaylistTracksOpt(
				subMUCoreSpotifyPlaylistID,
				&spotify.Options{Limit: &lim, Offset: &offset}, "",
			)
			if err != nil {
				if strings.Contains(err.Error(), "API rate limit exceeded") {
					return err
				}
				log.Printf("error fetching sub-/mu/core tracks from spotify")
				return nil
			}
			tt = append(tt, res.Tracks...)
			return nil
		}
		log.Println("fetching sub-/mu/core tracks from spotify...")
		exponentialBackoff(fn, 4)
	}

	var aa []Album
	for i := 0; i < len(tt); i++ {
		tr := tt[i]
		a := Album{
			Title:  tr.Track.Album.Name,
			Artist: tr.Track.Album.Artists[0].Name,
		}
		if imgs := tr.Track.Album.Images; len(imgs) > 0 {
			a.ArtworkURL = imgs[0].URL
		}
		aa = append(aa, a)
	}
	return uniqueAlbums(aa)
}

func uniqueAlbums(aa []Album) []Album {
	set := orderedset.New()
	for i := 0; i < len(aa); i++ {
		set.Add(aa[i])
	}
	var res []Album
	vv := set.Values()
	for i := 0; i < len(vv); i++ {
		res = append(res, vv[i].(Album))
	}

	return res
}
