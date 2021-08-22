package main

import (
	"github_ppva/do"
	"log"
)

func main() {
	if err := do.New().Run(); err != nil {
		log.Fatal(err)
	}
}
