package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"

	"github.com/kennygrant/sanitize"
	"github.com/nanoteck137/dwebble/collection"
	"github.com/nanoteck137/dwebble/utils"
)

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

var artistNames = []string{
	"Ado",
	"Metallica",
	"Bring Me The Horizon",
	"FZMZ",
	"Electric Callboy",
	"Slipknot",
	"Foo Fighters",
	"LiSA",
	"Survive Said The Prophet",
	"Eve",
	"SiM",
	"Vaundy",
	"YOASOBI",
	"System Of A Down",
	"In Flames",
	"Rammstein",
	"Linked Horizon",
	"Vickeblanka",
	"UVERworld",
	"SawanoHiroyuki[nZk]",
}

var albumNames = []string{
	"Ado's Utattemita Album",
	"ONE PIECE FILM RED",
	"Kyougen",
	"TEKKNO",
	"Metallica",
	"Master of Puppets",
	"We Are Not Your Kind",
	".5: The Grey Chapter",
	"All Hope Is Gone",
	"Vol. 3 The Subliminal Verses",
	"All Hope Is Gone",
	"Iowa",
	"Slipknot",
	"High Voltage",
	"Highway to Hell",
	"Back In Black",
}

var trackNames = []string{
	"DArkSide",
	"Can You Fell My Heart",
	"Throne",
	"Idol",
	"Mephisto",
	"Bokurano",
	"Velonica",
	"WORK",
	"Inferno",
	"KICK BACK",
	"The Rumbling",
	"Unsainted",
	"Nero Forte",
	"Critical Darling",
	"A Liar's Funeral",
	"unravel",
	"Stairway To Heaven",
	"Wash It All Away",
	"Song of the Dead",
}

func randomArtistName() string {
	index := randRange(0, len(artistNames))
	return artistNames[index]
}

func randomAlbumName() string {
	index := randRange(0, len(albumNames))
	return albumNames[index]
}

func randomTrackName() string {
	index := randRange(0, len(trackNames))
	return trackNames[index]
}

func getDir() (string, error) {
	dir := "./mockdata"
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
	} else {
		err := os.RemoveAll(dir)
		if err != nil {
			return "", err
		}
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func main() {
	// numArtists := randRange(2, 5)
	numArtists := 4

	var artists []collection.ArtistMetadata

	dir, err := getDir()
	if err != nil {
		log.Fatal(err)
	}

	// ffmpeg -f lavfi -i "sine=frequency=1000:duration=5" test.flac

	out := path.Join(dir, "testaudio.flac")
	utils.RunFFmpeg(false, "-f", "lavfi", "-i", "sine=frequency=1000:duration=5", out)

	for artistIndex := 0; artistIndex < numArtists; artistIndex++ {
		artistName := randomArtistName()
		fmt.Println("Artist", artistName)

		var albums []collection.AlbumMetadata

		for albumIndex := 0; albumIndex < 2; albumIndex++ {
			var tracks []collection.TrackMetadata

			for trackIndex := 0; trackIndex < 4; trackIndex++ {
				trackName := randomTrackName()

				tracks = append(tracks, collection.TrackMetadata{
					Id:       utils.CreateId(),
					Name:     trackName,
					FilePath: "",
					// Artist:   artistName,
				})
			}

			albumName := randomAlbumName()
			albums = append(albums, collection.AlbumMetadata{
				Id:     utils.CreateId(),
				Name:   albumName,
				Tracks: tracks,
			})
		}

		artists = append(artists, collection.ArtistMetadata{
			Id:     utils.CreateId(),
			Name:   artistName,
			Albums: albums,
		})
	}

	for _, artist := range artists {
		name := sanitize.Name(artist.Name)
		p := path.Join(dir, name)
		fmt.Printf("Path: %v\n", p)

		err := os.Mkdir(p, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}
