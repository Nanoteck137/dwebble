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
	"strconv"
	"strings"
	"unicode"

	"github.com/gosimple/slug"
	"github.com/mitchellh/mapstructure"
	"github.com/nanoteck137/parasect"
	"github.com/nrednav/cuid2"
)

func ExtractNumber(s string) int {
	n := ""
	for _, c := range s {
		if unicode.IsDigit(c) {
			n += string(c)
		} else {
			break
		}
	}

	if len(n) == 0 {
		return 0
	}

	i, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return 0
	}

	return int(i)
}

var CreateId = createIdGenerator(32)
var CreateSmallId = createIdGenerator(8)

func createIdGenerator(length int) func() string {
	res, err := cuid2.Init(cuid2.WithLength(length))
	if err != nil {
		log.Fatal("Failed to create id generator", "err", err)
	}

	return res
}

func Slug(s string) string {
	return slug.Make(s)
}

func SplitString(s string) []string {
	tags := []string{}
	if s != "" {
		tags = strings.Split(s, ",")
	}

	return tags
}

func Decode(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
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

var lossyFormatExts = []string{
	"opus",
	"mp3",
}

func IsLossyFormatExt(ext string) bool {
	return parasect.IsValidExt(lossyFormatExts, ext)
}

func IsFileLossyFormat(file string) bool {
	return IsLossyFormatExt(path.Ext(file))
}

func FormatTime(t int) string {
	s := t % 60
	m := t / 60

	return fmt.Sprintf("%v:%v", m, s)
}

func ParseAuthHeader(authHeader string) string {
	splits := strings.Split(authHeader, " ")
	if len(splits) != 2 {
		return ""
	}

	if splits[0] != "Bearer" {
		return ""
	}

	return splits[1]
}

func RunFFprobe(verbose bool, args ...string) ([]byte, error) {
	cmd := exec.Command("ffprobe", args...)
	if verbose {
		cmd.Stderr = os.Stderr
	}

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return data, nil
}

type TrackInfo struct {
	Path string

	Duration int
	Tags     map[string]string
}

type probeFormat struct {
	BitRate    string            `json:"bit_rate"`
	Tags       map[string]string `json:"tags"`
	Duration   string            `json:"duration"`
	FormatName string            `json:"format_name"`
}

type probeStream struct {
	Index     int    `json:"index"`
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`

	Duration string `json:"duration"`

	Tags map[string]string `json:"tags"`
}

type probe struct {
	Streams []probeStream `json:"streams"`
	Format  probeFormat   `json:"format"`
}

// TODO(patrik): Test
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

// TODO(patrik): Test
func convertMapKeysToLowercase(m map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

type ProbeResult struct {
	Tags     map[string]string
	Duration int
}

func ProbeTrack(filepath string) (ProbeResult, error) {
	data, err := RunFFprobe(false, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filepath)
	if err != nil {
		return ProbeResult{}, err
	}

	var probe probe
	err = json.Unmarshal(data, &probe)
	if err != nil {
		return ProbeResult{}, err
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
				return ProbeResult{}, err
			}

			duration = int(dur)
			if !hasGlobalTags {
				tags = convertMapKeysToLowercase(s.Tags)
			}
		}
	}

	return ProbeResult{
		Tags:     tags,
		Duration: duration,
	}, nil
}

func GetTrackInfo(filepath string) (TrackInfo, error) {
	probeResult, err := ProbeTrack(filepath)
	if err != nil {
		return TrackInfo{}, err
	}

	return TrackInfo{
		Path:     filepath,
		Duration: probeResult.Duration,
		Tags:     probeResult.Tags,
	}, nil
}

func CreateResizedImage(src string, dest string, dim int) error {
	args := []string{
		"convert",
		src,
		"-gravity", "Center",
		"-extent", "1:1",
		"-resize", fmt.Sprintf("%dx%d", dim, dim),
		dest,
	}

	cmd := exec.Command("magick", args...)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
