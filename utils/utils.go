package utils

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strconv"

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

type FileResult struct {
	Path   string
	Number int
	Name   string
}

var test1 = regexp.MustCompile(`(^\d+)[-\s]*(.+)`)
var test2 = regexp.MustCompile(`track(\d+).+`)

func CheckFile(filepath string) (FileResult, error) {
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

		return FileResult {
			Path: filepath,
			Number: num,
			Name: "",
		}, nil
	} else {
		num, err := strconv.Atoi(string(res[1]))
		if err != nil {
			return FileResult{}, nil
		}

		name := string(res[2])
		return FileResult {
			Path: filepath,
			Number: num,
			Name: name,
		}, nil
	}
}
