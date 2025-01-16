package utils

import (
	"os/exec"
	"path"
)

// TODO(patrik): Move this
func ProcessLossyVersion(input string, outputDir, name string) (string, error) {
	inputExt := path.Ext(input)
	isLossyInput := inputExt == ".opus" || inputExt == ".mp3"

	outputExt := ".opus"
	if isLossyInput {
		outputExt = inputExt
	}

	filename := name + outputExt

	var args []string
	// TODO(patrik): Set bitrate
	args = append(args, "-y", "-i", input, "-map_metadata", "-1", "-b:a", "128K")

	if outputExt == inputExt {
		args = append(args, "-codec", "copy")
	}

	args = append(args, path.Join(outputDir, filename))

	cmd := exec.Command("ffmpeg", args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return filename, nil
}

func ProcessOriginalVersion(input string, outputDir, name string) (string, error) {
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

	filename := name + outputExt

	var args []string
	args = append(args, "-i", input, "-map_metadata", "-1", "-vn")

	if outputExt == inputExt {
		args = append(args, "-codec", "copy")
	}

	out := path.Join(outputDir, filename)
	args = append(args, out)

	cmd := exec.Command("ffmpeg", args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return filename, nil
}
