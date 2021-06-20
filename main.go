package main

import (
	"github_ppva/do"
	"log"
)

func main() {
	if e := do.New().Run(); e != nil {
		log.Fatal(e)
	}
}
