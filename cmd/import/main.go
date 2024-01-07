package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/nanoteck137/dwebble/v2/collection"
	"github.com/nanoteck137/dwebble/v2/utils"
	"github.com/nrednav/cuid2"
	"github.com/spf13/cobra"
)

var validExts []string = []string{
	"wav",
	"m4a",
	"flac",
	"mp3",
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

func getOrCreateArtist(col *collection.Collection, metadata *metadata) *collection.ArtistDef {
	artistName := metadata.ArtistCredit[0].Name

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

	return artist
}

func getOrCreateAlbum(artist *collection.ArtistDef, metadata *metadata) *collection.AlbumDef {
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

	return album
}

func runImport(col *collection.Collection, importPath string) {
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
			fmt.Printf("p: %v\n", p)
			file, err := utils.CheckFile(p)
			if err == nil {
				fmt.Printf("p: %v %+v\n", p, file)
				files = append(files, file)
			} else {
				fmt.Printf("%v\n", err)
			}
		}
	}

	if len(files) == 0 {
		log.Fatal("No music files found")
	}

	fmt.Print("Enter MBID (musicbrainz.org): ")
	// mbid := "2529f558-970b-33d2-a42c-41ab15a970c6"
	// mbid := "eac4cbe6-00fa-4e3f-8304-e6674a34543d"

	reader := bufio.NewReader(os.Stdin)
	mbid, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

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

	if len(files) < len(m.Tracks) {
		fmt.Println("Warning: Missing tracks from album")
	}

	if len(files) > len(m.Tracks) {
		fmt.Println("Warning: More tracks then metadata album")
	}

	type res struct {
		file  utils.FileResult
		track track
	}

	mapping := []res{}

	fmt.Printf("Mapping:\n")
	for _, track := range m.Tracks {
		found := false
		for _, file := range files {
			if file.Number == track.Position {
				fmt.Printf("  %02v:%v -> %02v:%v:%v\n", track.Position, track.Title, file.Number, path.Base(file.Path), file.Name)
				found = true

				mapping = append(mapping, res{
					file:  file,
					track: track,
				})

				break
			}
		}

		if found {
			continue
		}

		fmt.Printf("  %02v:%v -> Missing\n", track.Position, track.Title)
	}

	artist := getOrCreateArtist(col, &metadata)
	album := getOrCreateAlbum(artist, &metadata)

	artistDirPath, err := col.GetArtistDir(artist)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("artistDirPath: %v\n", artistDirPath)

	albumPath := path.Join(artistDirPath, slug.Make(album.Name))
	err = os.MkdirAll(albumPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("albumPath: %v\n", albumPath)

	if len(album.Tracks) == 0 {
		// TODO(patrik): Add all
		for _, m := range mapping {
			id := gen()

			files, err := processFile(albumPath, id, m.file)
			if err != nil {
				log.Fatal(err)
			}

			prefix := path.Base(albumPath)
			files.Best = path.Join(prefix, files.Best)
			files.Mobile = path.Join(prefix, files.Mobile)
			fmt.Printf("files: %v\n", files)

			album.Tracks = append(album.Tracks, collection.TrackDef{
				Id:       id,
				Number:   uint(m.track.Position),
				CoverArt: "",
				Name:     m.track.Title,
				Files:    files,
			})
		}
	} else if len(album.Tracks) != len(mapping) {
		// TODO(patrik): Analyze missing tracks
		log.Printf("Mismatch\n")
	} else if len(album.Tracks) == len(mapping) {
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
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

var verbose = false

func runFFmpeg(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func processFile(outputDir string, id string, file utils.FileResult) (collection.TrackFiles, error) {
	// best - flac / maybe mp3
	// mobile - mp3

	fmt.Printf("m.file.Path: %v\n", file.Path)

	fileExt := path.Ext(file.Path)[1:]

	bestExt := fileExt
	// Set best extention to flac if the input is a wav file
	if fileExt == "wav" {
		bestExt = "flac"
	}

	mobileExt := "mp3"

	// Skip converting to mp3 if the input already is an mp3 file
	convertToMp3 := fileExt != "mp3"

	best := fmt.Sprintf("%v.%v.%v", id, "best", bestExt)
	mobile := best
	if convertToMp3 {
		mobile = fmt.Sprintf("%v.%v.%v", id, "mobile", mobileExt)
	}

	bestFilePath := path.Join(outputDir, best)

	// Convert wav files to flac
	if fileExt == "wav" {
		runFFmpeg("-y", "-i", file.Path, bestFilePath)
	} else {
		_, err := copy(file.Path, bestFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	if convertToMp3 {
		output := path.Join(outputDir, mobile)
		err := runFFmpeg("-y", "-i", file.Path, "-vn", "-ar", "44100", "-b:a", "192k", output)
		if err != nil {
			log.Fatal(err)
		}
	}

	return collection.TrackFiles{
		Best:   best,
		Mobile: mobile,
	}, nil
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

			runImport(col, importPath)

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
