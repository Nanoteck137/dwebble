package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/nrednav/cuid2"
)

var CreateId = createIdGenerator()

func createIdGenerator() func() string {
	res, err := cuid2.Init(cuid2.WithLength(32))
	if err != nil {
		log.Fatal(err)
	}

	return res
}


var validExts []string = []string{
	"wav",
	"m4a",
	"flac",
	"mp3",
}

func IsValidTrackExt(ext string) bool {
	for _, valid := range validExts {
		if valid == ext {
			return true
		}
	}

	return false
}

func CopyFile(src, dst string) (int64, error) {
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

func WriteFormFile(file *multipart.FileHeader, dst string) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	df, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, f)
	if err != nil {
		return err
	}

	return nil
}

func GetSingleFile(form *multipart.Form, name string) (*multipart.FileHeader, error) {
	files := form.File[name]
	if len(files) == 1 {
		return files[0], nil
	} else if len(files) > 1 {
		return nil, errors.New("Too many files, expected one file")
	}

	return nil, errors.New("Missing file")
}

var validMusicExts = []string{
	"flac",
	"opus",
	"mp3",
}

func IsMusicFile(p string) bool {
	if path.Base(p)[0] == '.' {
		return false
	}

	ext := path.Ext(p)
	if ext == "" {
		return false
	}

	ext = strings.ToLower(ext[1:])

	for _, validExt := range validMusicExts {
		if ext == validExt {
			return true
		}
	}

	return false
}

// TODO(patrik): Cleanup Multidisc stuff
var validMultiDiscPrefix = []string{
	"cd",
	"disc",
}

func IsMultiDiscName(name string) bool {
	for _, prefix := range validMultiDiscPrefix {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

func HasMusic(entries []fs.DirEntry) bool {
	for _, entry := range entries {
		if IsMusicFile(entry.Name()) {
			return true
		}
	}

	return false
}

func IsMultiDisc(entries []fs.DirEntry) bool {
	for _, entry := range entries {
		if IsMultiDiscName(entry.Name()) {
			return true
		}
	}
	return false
}

func SymlinkReplace(src, dst string) error {
	err := os.Symlink(src, dst)
	if err != nil {
		if os.IsExist(err) {
			err := os.Remove(dst)
			if err != nil {
				return err
			}

			err = os.Symlink(src, dst)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

var validImageExts = []string{
	"png",
	"jpg",
	"jpeg",
}

func IsValidImageExt(ext string) bool {
	for _, e := range validImageExts {
		if ext == e {
			return true
		}
	}

	return false
}

func FormatTime(t int) string {
	s := t % 60
	m := t / 60

	return fmt.Sprintf("%v:%v", m, s)
}
