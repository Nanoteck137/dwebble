package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
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

func RunFFprobe(args ...string) ([]byte, error) {
	cmd := exec.Command("ffprobe", args...)
	if true {
		cmd.Stderr = os.Stderr
	}

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func RunFFmpeg(verbose bool, args ...string) error {
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

type ProbeResult struct {
	Artist      string
	AlbumArtist string
	Title       string
	Album       string
	Track       int
	Disc        int
}

type FileResult struct {
	Path   string
	Number int
	Name   string

	Duration int
	Tags     map[string]string
}

type probeFormat struct {
	BitRate    string            `json:"bit_rate"`
	Tags       map[string]string `json:"tags"`
	Duration   string            `json:"duration"`
	FormatName string            `json:"format_name"`

	// "filename": "/Volumes/media/music/Various Artists/Cyberpunk 2077/cd1/19 - P.T. Adamczyk - Rite Of Passage.mp3",
	// "nb_streams": 2,
	// "nb_programs": 0,
	// "format_name": "mp3",
	// "format_long_name": "MP2/3 (MPEG audio layer 2/3)",
	// "start_time": "0.025056",
	// "size": "13898147",
	// "probe_score": 51,
}

type probeStream struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`

	Duration string `json:"duration"`

	// Video
	Width  int `json:"width"`
	Height int `json:"height"`

	Disposition struct {
		AttachedPic int `json:"attached_pic"`
		// "default": 0,
		// "dub": 0,
		// "original": 0,
		// "comment": 0,
		// "lyrics": 0,
		// "karaoke": 0,
		// "forced": 0,
		// "hearing_impaired": 0,
		// "visual_impaired": 0,
		// "clean_effects": 0,
		// "timed_thumbnails": 0,
		// "captions": 0,
		// "descriptions": 0,
		// "metadata": 0,
		// "dependent": 0,
		// "still_image": 0
	} `json:"disposition"`

	Tags map[string]string `json:"tags"`

	// "codec_long_name": "PNG (Portable Network Graphics) image",
	// "codec_tag_string": "[0][0][0][0]",
	// "codec_tag": "0x0000",
	// "coded_width": 512,
	// "coded_height": 512,
	// "closed_captions": 0,
	// "film_grain": 0,
	// "has_b_frames": 0,
	// "pix_fmt": "rgba",
	// "level": -99,
	// "color_range": "pc",
	// "refs": 1,
	// "r_frame_rate": "90000/1",
	// "avg_frame_rate": "0/0",
	// "time_base": "1/90000",
	// "start_pts": 2255,
	// "start_time": "0.025056",
	// "duration_ts": 30142433,
	// "duration": "334.915922",
}

type probe struct {
	Streams []probeStream `json:"streams"`
	Format  probeFormat   `json:"format"`
}

func getNumberFromFormatString(s string) int {
	if strings.Contains(s, "/") {
		s = strings.Split(s, "/")[0]
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}

	return num
}

func convertMapKeysToLowercase(m map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

var test1 = regexp.MustCompile(`(^\d+)[-\s]*(.+)\.`)
var test2 = regexp.MustCompile(`track(\d+).+`)

type Info struct {
	Tags     map[string]string
	Duration int
}

func GetInfo(filepath string) (Info, error) {
	// ffprobe -v quiet -print_format json -show_format -show_streams input
	data, err := RunFFprobe("-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filepath)
	if err != nil {
		return Info{}, err
	}

	var probe probe
	err = json.Unmarshal(data, &probe)
	if err != nil {
		return Info{}, err
	}

	hasGlobalTags := probe.Format.FormatName != "ogg"

	var tags map[string]string

	if hasGlobalTags {
		tags = convertMapKeysToLowercase(probe.Format.Tags)
	}

	duration := 0
	for _, s := range probe.Streams {
		if s.CodecType == "audio" {
			dur, err := strconv.ParseFloat(s.Duration, 32)
			if err != nil {
				return Info{}, err
			}

			duration = int(dur)
			if !hasGlobalTags {
				tags = convertMapKeysToLowercase(s.Tags)
			}
		}
	}

	return Info{
		Tags:     tags,
		Duration: duration,
	}, nil
}

func CheckFile(filepath string) (FileResult, error) {
	info, err := GetInfo(filepath)
	if err != nil {
		return FileResult{}, err
	}

	name := path.Base(filepath)
	res := test1.FindStringSubmatch(name)
	if res == nil {
		res := test2.FindStringSubmatch(name)
		if res == nil {
			return FileResult{}, fmt.Errorf("No result")
		}

		num, err := strconv.Atoi(string(res[1]))
		if err != nil {
			return FileResult{}, nil
		}

		return FileResult{
			Path:     filepath,
			Number:   num,
			Name:     "",
			Duration: info.Duration,
			Tags:     info.Tags,
		}, nil
	} else {
		num, err := strconv.Atoi(string(res[1]))
		if err != nil {
			return FileResult{}, nil
		}

		name := string(res[2])
		return FileResult{
			Path:     filepath,
			Number:   num,
			Name:     name,
			Duration: info.Duration,
			Tags:     info.Tags,
		}, nil
	}
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
