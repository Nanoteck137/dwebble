package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nrednav/cuid2"
)

func main() {
	length := flag.Int("length", 15, "Id length")
	flag.Parse()

	gen, err := cuid2.Init(cuid2.WithLength(*length))
	if err != nil {
		log.Fatal(err)
	}

	id := gen()
	fmt.Printf("%v\n", id)
}
