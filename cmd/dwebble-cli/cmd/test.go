package cmd

import (
	"context"
	"math"
	"strconv"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/spf13/cobra"
	"gopkg.in/vansante/go-ffprobe.v2"
)

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		p := "../audiotest/opus/01. test.opus"
		// p := "../audiotest/vorbis/01. test.ogg"
		// p := "../audiotest/flac/01. test.flac"
		// p := "../audiotest/aac/01. test.aac"
		// p := "../audiotest/01 The Rumbling.m4a"


		res, err := ffprobe.ProbeURL(context.TODO(), p)
		if err != nil {
			panic(err)
		}

		switch res.Format.FormatName {
		case "flac":
			log.Info("Found 'flac' audio")
		case "ogg":
			audioStream := res.FirstAudioStream()
			log.Info("Found 'ogg' container", "codec", audioStream.CodecName)
			s, _ := strconv.ParseFloat(audioStream.Duration, 32)
			s = math.Round(s)
			log.Info("Duration", "ts", int64(s))
		}
		
		pretty.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
