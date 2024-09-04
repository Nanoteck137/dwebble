package utils

import (
	"os"
	"os/exec"
	"path"
)

// TODO(patrik): Move this
func ProcessMobileVersion(input string, outputDir, name string) (string, error) {
	inputExt := path.Ext(input)
	isLossyInput := inputExt == ".opus" || inputExt == ".mp3"

	outputExt := ".opus"
	if isLossyInput {
		outputExt = inputExt
	}

	filename := Slug(name) + outputExt

	var args []string
	args = append(args, "-i", input, "-map_metadata", "-1")

	if outputExt == inputExt {
		args = append(args, "-codec", "copy")
	}

	args = append(args, path.Join(outputDir, filename))

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return filename, nil
}

func ProcessOriginalVersion(input string, outputDir, name string) (string, TrackInfo, error) {
	inputExt := path.Ext(input)
	// isLossyInput := inputExt == ".opus"

	outputExt := inputExt
	switch inputExt {
	case ".wav":
		outputExt = ".flac"
	}

	// .opus -> .opus
	// .mp3  -> .mp3
	// .wav  -> .flac
	// .flac -> .flac

	filename := Slug(name) + outputExt

	var args []string
	args = append(args, "-i", input, "-map_metadata", "-1")

	if outputExt == inputExt {
		args = append(args, "-codec", "copy")
	}

	out := path.Join(outputDir, filename)
	args = append(args, out)

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", TrackInfo{}, err
	}

	info, err := GetTrackInfo(input)
	if err != nil {
		return "", TrackInfo{}, err
	}

	return filename, info, nil
}
