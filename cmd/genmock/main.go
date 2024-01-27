package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/kennygrant/sanitize"
	"github.com/kr/pretty"
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

func fetchPlaceholderImage(dst string) error {
	url := "https://placehold.co/500x500.png"
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func generateRaw(dir string, audioFile string, imageFile string) {
	numArtists := 4

	for artistIndex := 0; artistIndex < numArtists; artistIndex++ {
		artistName := randomArtistName()
		artistName += " #" + strconv.Itoa(artistIndex)
		fmt.Println("Artist", artistName)

		name := sanitize.Name(artistName)
		artistDir := path.Join(dir, name)

		err := os.Mkdir(artistDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		audioFilePath, err := filepath.Abs(audioFile)
		if err != nil {
			log.Fatal(err)
		}

		imageFilePath, err := filepath.Abs(imageFile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v\n", audioFilePath)

		var albums []collection.AlbumMetadata
		for albumIndex := 0; albumIndex < 2; albumIndex++ {
			albumName := randomAlbumName()
			albumName += " #" + strconv.Itoa(albumIndex)

			name := sanitize.Name(albumName)
			albumDir := path.Join(artistDir, name)

			err := os.Mkdir(albumDir, 0755)
			if err != nil {
				log.Fatal(err)
			}

			var tracks []collection.TrackMetadata
			for trackIndex := 0; trackIndex < 4; trackIndex++ {
				trackName := randomTrackName()

				name := fmt.Sprintf("%v.flac", strconv.Itoa(trackIndex+1))
				os.Symlink(audioFilePath, path.Join(albumDir, name))

				tracks = append(tracks, collection.TrackMetadata{
					Id:       utils.CreateId(),
					Name:     trackName,
					Number:   trackIndex + 1,
					FileName: name,
					// ArtistId: "",
				})
			}

			os.Symlink(imageFilePath, path.Join(albumDir, "picture.png"))

			albums = append(albums, collection.AlbumMetadata{
				Id:       utils.CreateId(),
				Name:     albumName,
				CoverArt: "picture.png",
				Dir:      name,
				Tracks:   tracks,
			})
		}

		os.Symlink(imageFilePath, path.Join(artistDir, "picture.png"))

		artist := collection.ArtistMetadata{
			Id:      utils.CreateId(),
			Name:    artistName,
			Picture: "picture.png",
			Albums:  albums,
		}

		data, err := json.MarshalIndent(artist, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(path.Join(artistDir, "artist.json"), data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	// numArtists := randRange(2, 5)

	dir, err := getDir()
	if err != nil {
		log.Fatal(err)
	}

	// ffmpeg -f lavfi -i "sine=frequency=1000:duration=5" test.flac

	audioFile := path.Join(dir, "testaudio.flac")
	utils.RunFFmpeg(false, "-f", "lavfi", "-i", "sine=frequency=1000:duration=5", audioFile)

	imageFile := path.Join(dir, "image.png")
	_, err = os.Stat(imageFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = fetchPlaceholderImage(imageFile)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	image, err := os.Open(imageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer image.Close()

	col := collection.NewEmpty(dir)
	artist, err := col.CreateArtist("Test Artist", collection.File{
		Content:     image,
		ContentType: "image/png",
	})
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(artist)

	album, err := col.CreateAlbum("Test Album", collection.File{
		Content:     image,
		ContentType: "image/png",
	}, artist)
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(album)
	pretty.Println(col)

	err = col.FlushToDisk()
	if err != nil {
		log.Fatal(err)
	}
}
