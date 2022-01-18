package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {

	var (
		chastity        = flag.Float64("chastity", 0.5, "Chastity filter")
		merLen          = flag.Int("merlen", 8, "Mer length")
		consideredBases = flag.Int("consideredbases", 100, "Number of bases to consider")
		file            = flag.String("file", "", "File to process")
	)

	flag.Parse()

	if *file == "" {
		fmt.Println("Usage: entropy-calc -file <filename> -chastity <float> -merlen <int> -consideredbases <int>")
		os.Exit(1)
	}

	// open examples/input.reads
	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	matchedMers := CreateMatchedMers(*chastity, *consideredBases)

	// create a channel to read from

	coll := CreateMerMatchCollection(*chastity, *merLen, *consideredBases, &matchedMers)

	// start waiting for lines
	coll.Start()

	// read from the file and send to the channel
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		coll.Broadcast(scanner.Text())
	}

	// close the channel
	coll.Done()

	// wait for the goroutines to finish
	coll.Wait()

	// summarise the results
	matchedMers.Summarise()

	// close the file
	f.Close()

}
