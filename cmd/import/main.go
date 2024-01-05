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
	"github.com/nrednav/cuid2"
	"github.com/spf13/cobra"
	"github.com/gosimple/slug"
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

			return

			entries, err := os.ReadDir(importPath)
			if err != nil {
				log.Fatal(err)
			}

			// TODO(patrik): Check if sorting of files is working correctly
			for _, entry := range entries {
				p := path.Join(importPath, entry.Name())
				ext := path.Ext(p)[1:]

				if isValidExt(ext) {
					fmt.Printf("p: %v\n", p)
				}
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

			tracks := make([]collection.TrackDef, len(m.Tracks))

			for i, track := range m.Tracks {
				tracks[i] = collection.TrackDef{
					Id:       gen(),
					Number:   uint(track.Position),
					CoverArt: "",
					Name:     track.Title,
				}
			}

			albumDef := collection.AlbumDef{
				Id:       gen(),
				Name:     metadata.Title,
				CoverArt: "",
				Tracks:   tracks,
			}

			artistDef := collection.ArtistDef{
				Id:      gen(),
				Name:    metadata.ArtistCredit[0].Name,
				Changed: 0,
				Picture: "",
				Albums: []collection.AlbumDef{
					albumDef,
				},
			}

			j, err := json.MarshalIndent(artistDef, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(j))

			s := slug.Make(artistDef.Name)
			fmt.Printf("s: %v\n", s)

			p := path.Join(collectionPath, s)
			err = os.MkdirAll(p, 0755)
			if err != nil {
				log.Fatal(err)
			}

			// err = os.WriteFile(path.Join(p, "artist.json"), j, 0644)
			// if err != nil {
			// 	log.Fatal(err)
			// }
		},
	}

	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
