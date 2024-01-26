package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nanoteck137/dwebble/collection"
	"github.com/nanoteck137/dwebble/utils"
)

func main() {
	col, err := collection.ReadFromDir("../testcol")
	if err != nil {
		log.Fatal(err)
	}

	_ = col

	return

	// otherartist := collection.ArtistMetadata{
	// 	Id:     utils.CreateId(),
	// 	Name:   "Some Other Artist",
	// 	Albums: []collection.AlbumMetadata{},
	// }

	artist := collection.ArtistMetadata{
		Id:   utils.CreateId(),
		Name: "Ado",
		Albums: []collection.AlbumMetadata{
			{
				Id:   utils.CreateId(),
				Name: "unravel",
				Tracks: []collection.TrackMetadata{
					{
						Id:       utils.CreateId(),
						Name:     "unravel",
						FilePath: "unravel/01 - unravel.flac",
					},
				},
			},
			{
				Id:   utils.CreateId(),
				Name: "Kura Kura",
				Tracks: []collection.TrackMetadata{
					{
						Id:       utils.CreateId(),
						Name:     "Kura Kura",
						FilePath: "Kura Kura/01. クラクラ.flac",
					},
				},
			},
		},
	}

	data, err := json.MarshalIndent(artist, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))

	err = os.WriteFile("../testcol/Ado/artist.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// col := collection.New()
	//
	// data, err := json.Marshal(col)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println(string(data))
}
