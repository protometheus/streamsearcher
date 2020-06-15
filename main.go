package main

import (
	"log"
	"os"
)

func main() {
	os.Remove("./_input.txt")
	GenerateInput(1000000000)

	file, err := os.Open("./_input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("file size: %v\n", fs.Size())

	streamSearcher, err := NewStreamSearcher()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Stream Searcher...")
	streamSearcher.Run()
}
