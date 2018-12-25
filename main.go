package main

import (
	"log"
	"os"

	"./ui"
)

func main() {
	log.SetFlags(0)
	ui.Run(os.Args[1])
}
