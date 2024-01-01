// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package database

import ()

type Album struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistID string `json:"artistId"`
}

type Artist struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Track struct {
	ID          string `json:"id"`
	TrackNumber int32  `json:"trackNumber"`
	Name        string `json:"name"`
	CoverArt    string `json:"coverArt"`
	Filename    string `json:"filename"`
	AlbumID     string `json:"albumId"`
	ArtistID    string `json:"artistId"`
}
