package imp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/v2/collection"
	"github.com/nanoteck137/dwebble/v2/musicbrainz"
	"github.com/nanoteck137/dwebble/v2/utils"
)

type ImageFile struct {
	File   string
	Remove bool
}

type ImportMetadataTrack struct {
	Path   string
	Title  string
	Number int
}

type ImportMetadata struct {
	AlbumName       string
	AlbumArtistName string
	CoverArt        ImageFile
	Tracks          []ImportMetadataTrack
}

func getOrCreateArtist(col *collection.Collection, metadata *ImportMetadata) *collection.ArtistDef {
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

func getOrCreateAlbum(artist *collection.ArtistDef, metadata *ImportMetadata) *collection.AlbumDef {
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

func getFiles(p string) ([]utils.FileResult, error) {
	dirEntries, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var entries []utils.FileResult

	for _, entry := range dirEntries {
		fullPath := path.Join(p, entry.Name())

		if entry.IsDir() {
			continue
		}

		ext := path.Ext(fullPath)[1:]
		if utils.IsValidTrackExt(ext) {
			file, err := utils.CheckFile(fullPath)
			if err != nil {
				log.Printf("%v", err)
				continue
			}

			entries = append(entries, file)
		}
	}

	pretty.Println(entries)

	return entries, nil
}

const (
	FILE_STATUS_OK            = 0
	FILE_STATUS_MISSING_TITLE = 1 << 0
)

func checkFiles(files []utils.FileResult) bool {
	valid := true

	// fileStatus := make([]int, len(files))

	tracks := make(map[int][]int)

	for i, file := range files {
		if file.Name == "" && file.Probe.Title == "" {
			fmt.Printf("No title available\n")
			valid = false
		}

		if file.Probe.Artist == "" {
			fmt.Printf("No artist assigned\n")
			valid = false
		}

		if file.Probe.AlbumArtist == "" {
			fmt.Printf("No album artist assigned\n")
			valid = false
		}

		if file.Probe.Album == "" {
			fmt.Printf("No album assigned\n")
			valid = false
		}

		if file.Probe.Track == -1 {
			fmt.Printf("No track assigned\n")
			valid = false
		}

		if file.Probe.Disc == -1 {
			fmt.Printf("No disc assigned\n")
			valid = false
		}

		tracks[file.Number] = append(tracks[file.Number], i)

		// Artist      string
		// AlbumArtist string
		// Album       string
		// Track       int
		// Disc        int
	}

	for n, t := range tracks {
		if len(t) > 1 {
			log.Fatalf("Mutliple tracks with the number %v\n", n)
		}
	}

	return valid
}

func askForMbid() uuid.UUID {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter MBID: ")
	text, _ := reader.ReadString('\n')

	id, err := uuid.Parse(strings.TrimSpace(text))
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func missingMetadata(files []utils.FileResult, importMetadata *ImportMetadata) {
	mbid := askForMbid()
	fmt.Printf("MBID: %v\n", mbid)

	albumMetadata, err := musicbrainz.FetchAlbumMetadata(mbid.String())
	if err != nil {
		log.Fatal(err)
	}

	importMetadata.AlbumArtistName = albumMetadata.ArtistCredit[0].Name
	importMetadata.AlbumName = albumMetadata.Title

	type t struct {
		path   string
		title  string
		number int
	}

	media := albumMetadata.Media[0]
	var tracks []t

	for _, file := range files {
		var meta musicbrainz.Track
		found := false

		for _, tm := range media.Tracks {
			if tm.Position == file.Number {
				meta = tm
				found = true
			}
		}

		if !found {
			log.Fatalf("Failed to find mapping for %v\n", file.Path)
		}

		tracks = append(tracks, t{
			path:   file.Path,
			title:  meta.Title,
			number: file.Number,
		})
	}

	// TODO(patrik): Maybe sort here?

	for _, track := range tracks {
		importMetadata.Tracks = append(importMetadata.Tracks, ImportMetadataTrack{
			Path:   track.path,
			Title:  track.title,
			Number: track.number,
		})
	}

	// NOTE(patrik): Fetch the cover art from the archive
	coverArt, err := musicbrainz.FetchCoverArt(mbid.String())
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("*.%v", coverArt.Ext))
	defer file.Close()

	_, err = file.Write(coverArt.Data)
	if err != nil {
		log.Fatal(err)
	}

	importMetadata.CoverArt = ImageFile{
		File:   file.Name(),
		Remove: true,
	}

	fmt.Printf("Cover Art: %v\n", file.Name())
}

func ImportDir(col *collection.Collection, importPath string) error {
	files, err := getFiles(importPath)
	if err != nil {
		return err
	}

	importMetadata := ImportMetadata{}

	if checkFiles(files) {
		f := &files[0]
		importMetadata.AlbumArtistName = f.Probe.AlbumArtist
		importMetadata.AlbumName = f.Probe.Album

		for _, f := range files {
			importMetadata.Tracks = append(importMetadata.Tracks, ImportMetadataTrack{
				Path:   f.Path,
				Title:  f.Probe.Title,
				Number: f.Number,
			})
		}
	} else {
		missingMetadata(files, &importMetadata)
	}

	// if checkFiles() {
	// 	if askUserIfTheyWantToFetchMetadataFromDatabase() {
	// 		askUserForMBID()
	// 		fetchMetadataFromDatabase()
	// 		createImportMetadataFromDatabase()
	// 	} else {
	// 		createImportMetadataFromFiles()
	// 	}
	// } else {
	// 	if askUserIfTheyWantToManuallyEnterMetadata() {
	// 		askUserForData()
	// 		createImportMetadataFromUserInput()
	// 	} else {
	// 		askUserForMBID()
	// 		fetchMetadataFromDatabase()
	// 		createImportMetadataFromDatabase()
	// 	}
	// }

	pretty.Print(importMetadata)
	fmt.Println()

	artist := getOrCreateArtist(col, &importMetadata)
	album := getOrCreateAlbum(artist, &importMetadata)

	if importMetadata.CoverArt.File != "" {
		coverArtFile := importMetadata.CoverArt.File

		dir, err := col.GetArtistDir(artist)
		if err != nil {
			log.Fatal(err)
		}

		imagesDir := path.Join(dir, "images")
		err = os.MkdirAll(imagesDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		ext := path.Ext(coverArtFile)[1:]
		name := fmt.Sprintf("cover-art.original.%v", ext)

		_, err = utils.CopyFile(coverArtFile, path.Join(imagesDir, name))
		if err != nil {
			log.Fatal(err)
		}

		album.CoverArt = path.Join("images", name)

		if importMetadata.CoverArt.Remove {
			log.Printf("Removing cover art: %v\n", coverArtFile)
			err := os.Remove(coverArtFile)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if len(album.Tracks) == 0 {
		// TODO(patrik): Just add all tracks
		for _, track := range importMetadata.Tracks {
			id := utils.CreateId()
			// TODO(patrik): Process the file

			album.Tracks = append(album.Tracks, collection.TrackDef{
				Id:       id,
				Number:   uint(track.Number),
				CoverArt: "TODO",
				Name:     track.Title,
				Files: collection.TrackFiles{
					Best:   "TODO",
					Mobile: "TODO",
				},
			})
		}
	} else {
		// TODO(patrik): Check to see what needs to be done here
		log.Println("Check")
	}

	return nil
}
