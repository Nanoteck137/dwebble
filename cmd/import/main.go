package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/nanoteck137/dwebble/v2/collection"
	"github.com/nanoteck137/dwebble/v2/utils"
	"github.com/nrednav/cuid2"
	"github.com/spf13/cobra"
)

var validExts []string = []string{
	"wav",
}

func isValidExt(ext string) bool {
	for _, valid := range validExts {
		if valid == ext {
			return true
		}
	}

	return false
}

type track struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Number   string `json:"number"`
	Position int    `json:"position"`

	//	          "number": "2",
	//	          "length": 324600,
	//	          "position": 2,
	//	          "title": "Sad but True",

	Recording struct {
		ArtistCredit []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			//	                "joinphrase": "",
			//	                "name": "Metallica",
			//	                "artist": {
			//	                  "type": "Group",
			//	                  "sort-name": "Metallica",
			//	                  "type-id": "e431f5f6-b5d2-343d-8b36-72607fffb74b",
			//	                  "disambiguation": "",
			//	                  "name": "Metallica",
			//	                  "id": "65f4f0c5-ef9e-490c-aee3-909e7ae6b2ab"
			//	                }
		} `json:"artist-credit"`
		Disambiguation   string `json:"disambiguation"`
		Title            string `json:"title"`
		Length           int    `json:"length"`
		Id               string `json:"id"`
		FirstReleaseDate string `json:"first-release-date"`
		Video            bool   `json:"video"`
	} `json:"recording"`

	//	          "number": "2",
	//	          "length": 324600,
	//	          "position": 2,
	//	          "title": "Sad but True",
	//	          "artist-credit": [
	//	            {
	//	              "joinphrase": "",
	//	              "name": "Metallica",
	//	              "artist": {
	//	                "id": "65f4f0c5-ef9e-490c-aee3-909e7ae6b2ab",
	//	                "name": "Metallica",
	//	                "type-id": "e431f5f6-b5d2-343d-8b36-72607fffb74b",
	//	                "disambiguation": "",
	//	                "sort-name": "Metallica",
	//	                "type": "Group"
	//	              }
	//	            }
	//	          ]
}

type media struct {
	Title    string `json:"title"`
	FormatId string `json:"format-id"`
	Position int    `json:"position"`
	Format   string `json:"format"`

	TrackCount  string `json:"track-count"`
	TrackOffset string `json:"track-offset"`

	Tracks []track `json:"tracks"`
}

type metadata struct {
	Id    string  `json:"id"`
	Title string  `json:"title"`
	Date  string  `json:"date"`
	Media []media `json:"media"`

	ArtistCredit []struct {
		Name   string `json:"name"`
		Artist struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"arist"`
	} `json:"artist-credit"`

	//	{
	//	  "packaging-id": null,
	//	  "text-representation": {
	//	    "language": "eng",
	//	    "script": "Latn"
	//	  },
	//	  "title": "Metallica",
	//	  "cover-art-archive": {
	//	    "front": true,
	//	    "artwork": true,
	//	    "count": 6,
	//	    "back": true,
	//	    "darkened": false
	//	  },
	//	  "asin": "B000002H97",
	//	  "disambiguation": "",
	//	  "status": "Official",
	//	  "release-events": [
	//	    {
	//	      "date": "1991-08-12",
	//	      "area": {
	//	        "type-id": null,
	//	        "disambiguation": "",
	//	        "iso-3166-1-codes": [
	//	          "CA"
	//	        ],
	//	        "type": null,
	//	        "sort-name": "Canada",
	//	        "id": "71bbafaa-e825-3e15-8ca9-017dcad1748b",
	//	        "name": "Canada"
	//	      }
	//	    }
	//	  ],
	//	  "packaging": null,
	//	  "barcode": "075596111324",
	//	  "id": "2529f558-970b-33d2-a42c-41ab15a970c6",
	//	  "country": "CA",
	//	  "quality": "high",
	//	  "date": "1991-08-12",
	//	  "artist-credit": [
	//	    {
	//	    }
	//	  ],
	//	  "status-id": "4e304316-386d-3409-af2e-78857eec5cfe"
	//	}
}

func fetchAlbumMetadata(mbid string) (metadata, error) {
	// https://musicbrainz.org/ws/2/release/2529f558-970b-33d2-a42c-41ab15a970c6?inc=artist-credits%2Brecordings&fmt=json
	url := fmt.Sprintf("https://musicbrainz.org/ws/2/release/%v?inc=artist-credits+recordings&fmt=json", mbid)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return metadata{}, err
	}

	req.Header.Set("User-Agent", "dwebble/0.0.1 ( github.com/nanoteck137/dwebble )")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return metadata{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return metadata{}, err
	}

	var metadata metadata
	json.Unmarshal(data, &metadata)

	fmt.Printf("Title: %v\n", metadata.Title)
	fmt.Printf("Date: %v\n", metadata.Date)
	for _, m := range metadata.Media {
		fmt.Printf("  %v\n", m.Position)

		for _, t := range m.Tracks {
			fmt.Printf("    %v", t.Recording.Title)
			fmt.Printf(" - %v\n", t.Recording.ArtistCredit[0].Name)
		}
	}

	return metadata, nil
}

func main() {
	root := &cobra.Command{
		Use:   "dwebble-import",
		Short: "Import music into a collection that Dwebble can use",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("args: %v\n", args)

			collectionPath := args[0]
			importPath := args[1]

			fmt.Printf("collectionPath: %v\n", collectionPath)
			fmt.Printf("importPath: %v\n", importPath)

			col, err := collection.Read(collectionPath)
			if err != nil {
				log.Fatal(err)
			}

			_ = col

			entries, err := os.ReadDir(importPath)
			if err != nil {
				log.Fatal(err)
			}

			files := []utils.FileResult{}

			// TODO(patrik): Check if sorting of files is working correctly
			for _, entry := range entries {
				p := path.Join(importPath, entry.Name())
				ext := path.Ext(p)[1:]

				if isValidExt(ext) {
					file, err := utils.CheckFile(p)
					if err == nil {
						fmt.Printf("p: %v %+v\n", p, file)
						files = append(files, file)
					}
				}
			}

			if len(files) == 0 {
				log.Fatal("No music files found")
			}

			// reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter MBID (musicbrainz.org): ")
			mbid := "2529f558-970b-33d2-a42c-41ab15a970c6"
			// mbid, err := reader.ReadString('\n')
			// if err != nil {
			// 	log.Fatal(err)
			// }

			mbid = strings.TrimSpace(mbid)

			fmt.Println("MBID:", mbid)

			_, err = uuid.Parse(mbid)
			if err != nil {
				log.Fatal(err)
			}

			metadata, err := fetchAlbumMetadata(mbid)
			if err != nil {
				log.Fatal(err)
			}

			gen, err := cuid2.Init(cuid2.WithLength(32))
			if err != nil {
				log.Fatal(err)
			}

			m := metadata.Media[0]

			artistName := metadata.ArtistCredit[0].Name
			fmt.Printf("artistName: %v\n", artistName)

			if len(files) != len(m.Tracks) {
				fmt.Println("Warning: Missing tracks from album")
			}

			fmt.Printf("Mapping:\n")
			for i, track := range m.Tracks {
				if i < len(files) {
					fmt.Printf("  %02v:%v -> %02v:%v:%v\n", track.Position, track.Title, files[i].Number, path.Base(files[i].Path), files[i].Name)
				} else {
					fmt.Printf("  %02v:%v -> Missing\n", track.Position, track.Title)
				}
			}

			artist, err := col.FindArtistByName(artistName)
			if err != nil {
				if err == collection.NotFoundErr {
					log.Printf("Artist '%v' not found in collection (creating)", artistName)
					artist = col.CreateArtist(artistName)
				} else {
					log.Fatal(err)
				}
			} else {
				// TODO(patrik): Check artist against metadata
			}

			album, err := artist.FindAlbumByName(metadata.Title)
			if err != nil {
				if err == collection.NotFoundErr {
					album = artist.CreateAlbum(metadata.Title)
				} else {
					log.Fatal(err)
				}
			} else {
				// TODO(patrik): Check album against metadata
			}

			if len(album.Tracks) == 0 {
				// TODO(patrik): Add all
				for _, track := range m.Tracks {
					album.Tracks = append(album.Tracks, collection.TrackDef{
						Id:       gen(),
						Number:   uint(track.Position),
						CoverArt: "",
						Name:     track.Title,
					})
				}
			} else if len(album.Tracks) != len(m.Tracks) {
				// TODO(patrik): Analyze missing tracks
				log.Printf("Mismatch\n")
			} else if len(album.Tracks) == len(m.Tracks) {
				// TODO(patrik): Check so album.Tracks and m.Tracks matches
				log.Printf("Check\n")

				// for i := 0; i < len(album.Tracks); i++ {
				// 	mtrack := m.Tracks[i]
				// 	atrack := album.Tracks[i]
				//
				// 	if atrack.Name != mtrack.Title {
				// 		fmt.Printf("Album and Metadata track name not matching '%v' : '%v'\n", atrack.Name, mtrack.Title)
				// 	}
				// }
			} else {
				log.Fatal("No metadata tracks?")
			}

			err = col.Flush()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}