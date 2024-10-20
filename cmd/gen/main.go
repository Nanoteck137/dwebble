package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nanoteck137/dwebble/apis"
	"github.com/nanoteck137/pyrin/spec"
)

func main() {
	router := spec.Router{}

	apis.RegisterHandlers(nil, &router)

	s, err := spec.GenerateSpec(router.Routes)
	if err != nil {
		log.Fatal("Failed to generate spec", err)
	}

	d, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal server", err)
	}

	out := "./misc/pyrin.json"
	err = os.WriteFile(out, d, 0644)
	if err != nil {
		log.Fatal("Failed to write pyrin.json", err)
	}

	fmt.Println("Wrote 'misc/pyrin.json'")
}
