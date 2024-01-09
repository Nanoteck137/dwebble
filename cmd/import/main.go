package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"

	"github.com/gosimple/slug"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/v2/collection"
	"github.com/nanoteck137/dwebble/v2/imp"
	"github.com/nanoteck137/dwebble/v2/utils"
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

type metadataTrack struct {
	Path       string
	Title      string
	ArtistName string
	Number     int
}

type metadata struct {
	AlbumArtistName string
	AlbumName       string

	tracks []metadataTrack
}

func getOrCreateArtist(col *collection.Collection, metadata *metadata) *collection.ArtistDef {
	artistName := metadata.AlbumArtistName

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
	albumName := metadata.AlbumName

	album, err := artist.FindAlbumByName(metadata.AlbumName)
	if err != nil {
		if err == collection.NotFoundErr {
			album = artist.CreateAlbum(albumName)
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

	// Metadata
	// Artist
	// Album
	// All Tracks (Available, Missing)

	// Metadata
	//  Filename - Track, Maybe name
	//  Musicdb - All
	//  ffprobe -

	albumName := ""
	albums := make(map[string]int)
	for _, file := range files {
		if file.Probe.Album != "" {
			albumName = file.Probe.Album
			albums[file.Probe.Album]++
		}
	}

	if len(albums) == 0 {
		log.Fatal("No album name detected")
	}

	if len(albums) > 1 {
		log.Fatal("Multiple album name detected")
	}

	if albumName == "" {
		// TODO(patrik): Let the user provide this (MBID)
		log.Fatal("Should not happend")
	}

	fmt.Printf("albumName: %v\n", albumName)

	artists := make(map[string]int)
	albumArtist := ""

	for _, file := range files {
		if file.Probe.AlbumArtist != "" {
			albumArtist = file.Probe.AlbumArtist
		}

		if file.Probe.Artist != "" {
			artists[file.Probe.Artist]++
		}
	}

	if albumArtist == "" {
		if len(artists) > 1 {
			albumArtist = "Various Artists"
		} else {
			// TODO(patrik): Let the user provide this (MBID)
			log.Fatal("No album artist found")
		}
	}

	fmt.Printf("albumArtist: %v\n", albumArtist)

	const (
		STATUS_OK             = 0
		STATUS_MISSING_TITLE  = 1 << 0
		STATUS_MISSING_ARTIST = 1 << 1
		// TODO(patrik): Detect
		STATUS_MISSING_TRACK = 1 << 2
	)

	type t struct {
		Path   string
		Number int
		Title  string
		Artist string
		Status int
	}

	tracks := make(map[int]t)

	for _, file := range files {
		num := file.Number
		if file.Probe.Track != file.Number {
			log.Printf("Track number not matching (using filename number)")
		}

		if tracks[num].Path != "" {
			log.Fatal("Multiple tracks with the name number")
		}

		status := STATUS_OK
		title := file.Probe.Title
		artist := file.Probe.Artist

		if title == "" {
			status |= STATUS_MISSING_TITLE
		}

		if artist == "" {
			status |= STATUS_MISSING_ARTIST
		}

		tracks[num] = t{
			Path:   file.Path,
			Number: num,
			Title:  title,
			Artist: artist,
			Status: status,
		}
	}

	fmt.Printf("tracks: %+v\n", tracks)

	metadata := metadata{
		AlbumArtistName: albumArtist,
		AlbumName:       albumName,
	}

	for _, track := range tracks {
		title := track.Title
		// TODO(patrik): Handle this
		if track.Status&STATUS_MISSING_TITLE > 0 {
			title = "MISSING_TITLE"
		}

		artist := track.Artist
		// TODO(patrik): Handle this
		if track.Status&STATUS_MISSING_ARTIST > 0 {
			artist = "MISSING_ARTIST"
		}

		fmt.Printf("track %v:\n", track.Number)
		fmt.Printf("  title: %v\n", title)
		fmt.Printf("  artist: %v\n", artist)

		metadata.tracks = append(metadata.tracks, metadataTrack{
			Path:       track.Path,
			Title:      title,
			ArtistName: artist,
			Number:     track.Number,
		})
	}

	sort.SliceStable(metadata.tracks, func(i, j int) bool {
		return metadata.tracks[i].Number < metadata.tracks[j].Number
	})
	pretty.Print(metadata)

	// fmt.Print("Enter MBID (musicbrainz.org): ")
	// // mbid := "2529f558-970b-33d2-a42c-41ab15a970c6"
	// // mbid := "eac4cbe6-00fa-4e3f-8304-e6674a34543d"
	//
	// reader := bufio.NewReader(os.Stdin)
	// mbid, err := reader.ReadString('\n')
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// mbid = strings.TrimSpace(mbid)
	//
	// fmt.Println("MBID:", mbid)
	//
	// _, err = uuid.Parse(mbid)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// metadata, err := musicbrainz.FetchAlbumMetadata(mbid)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// gen, err := cuid2.Init(cuid2.WithLength(32))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// m := metadata.Media[0]
	//
	// if len(files) < len(m.Tracks) {
	// 	fmt.Println("Warning: Missing tracks from album")
	// }
	//
	// if len(files) > len(m.Tracks) {
	// 	fmt.Println("Warning: More tracks then metadata album")
	// }
	//
	// type res struct {
	// 	file  utils.FileResult
	// 	track track
	// }
	//
	// mapping := []res{}
	//
	// fmt.Printf("Mapping:\n")
	// for _, track := range m.Tracks {
	// 	found := false
	// 	for _, file := range files {
	// 		if file.Number == track.Position {
	// 			fmt.Printf("  %02v:%v -> %02v:%v:%v\n", track.Position, track.Title, file.Number, path.Base(file.Path), file.Name)
	// 			found = true
	//
	// 			mapping = append(mapping, res{
	// 				file:  file,
	// 				track: track,
	// 			})
	//
	// 			break
	// 		}
	// 	}
	//
	// 	if found {
	// 		continue
	// 	}
	//
	// 	fmt.Printf("  %02v:%v -> Missing\n", track.Position, track.Title)
	// }

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
		for _, track := range metadata.tracks {
			id := utils.CreateId()

			files, err := processFile(albumPath, id, track)
			if err != nil {
				log.Fatal(err)
			}

			prefix := path.Base(albumPath)
			files.Best = path.Join(prefix, files.Best)
			files.Mobile = path.Join(prefix, files.Mobile)
			fmt.Printf("files: %v\n", files)

			album.Tracks = append(album.Tracks, collection.TrackDef{
				Id:       id,
				Number:   uint(track.Number),
				CoverArt: "",
				Name:     track.Title,
				Files:    files,
			})
		}
	} else if len(album.Tracks) != len(metadata.tracks) {
		// TODO(patrik): Analyze missing tracks
		log.Printf("Mismatch\n")
	} else if len(album.Tracks) == len(metadata.tracks) {
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

func processFile(outputDir string, id string, track metadataTrack) (collection.TrackFiles, error) {
	// best - flac / maybe mp3
	// mobile - mp3

	fmt.Printf("m.file.Path: %v\n", track.Path)

	fileExt := path.Ext(track.Path)[1:]

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
		runFFmpeg("-y", "-i", track.Path, bestFilePath)
	} else {
		_, err := utils.CopyFile(track.Path, bestFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	if convertToMp3 {
		output := path.Join(outputDir, mobile)
		err := runFFmpeg("-y", "-i", track.Path, "-vn", "-ar", "44100", "-b:a", "192k", output)
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

			err = imp.ImportDir(col, importPath)
			if err != nil {
				log.Fatal(err)
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
