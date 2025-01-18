package utils

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func convertMapKeysToLowercase(m map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

type TrackInfo struct {
	Path string

	Duration int
	Tags     map[string]string
}

type ProbeResult struct {
	Tags     ffprobe.Tags
	Duration float64
}

func ProbeTrack(filepath string) (ProbeResult, error) {
	ctx := context.TODO()

	probe, err := ffprobe.ProbeURL(ctx, filepath)
	if err != nil {
		return ProbeResult{}, nil
	}

	var tags ffprobe.Tags
	hasGlobalTags := probe.Format.FormatName != "ogg"

	audioStream := probe.FirstAudioStream()
	if audioStream == nil {
		return ProbeResult{}, errors.New("contains no audio streams")
	}

	if hasGlobalTags {
		tags = probe.Format.TagList
	} else {
		tags = audioStream.TagList
	}

	tags = convertMapKeysToLowercase(tags)

	duration, err := strconv.ParseFloat(audioStream.Duration, 32)
	if err != nil {
		return ProbeResult{}, err
	}

	return ProbeResult{
		Tags:     tags,
		Duration: duration,
	}, nil
}
