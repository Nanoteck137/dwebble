package imp

import (
	"fmt"
	"log"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/v2/collection"
	"github.com/nanoteck137/dwebble/v2/utils"
)

type ImportMetadataTrack struct {
	Path   string
	Title  string
	Number int
}

type ImportMetadata struct {
	AlbumName       string
	AlbumArtistName string
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

type File struct {
	Path string
	Number int
}

func getFiles(p string) []File {
	return nil
}

func ImportDir(col *collection.Collection, importPath string) error {
	files := getFiles(importPath)
	_ = files

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

	importMetadata := ImportMetadata{
		AlbumName:       "Test Album",
		AlbumArtistName: "Test Artist",
		Tracks: []ImportMetadataTrack{
			{
				Path:   "/some/path",
				Title:  "Track #1",
				Number: 1,
			},
			{
				Path:   "/some/other/path",
				Title:  "Track #2",
				Number: 2,
			},
		},
	}

	pretty.Print(importMetadata)
	fmt.Println()

	artist := getOrCreateArtist(col, &importMetadata)
	album := getOrCreateAlbum(artist, &importMetadata)

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
