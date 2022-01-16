package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
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

	matchedMers := CreateMatchedMers(*chastity)
	staticLen := *consideredBases
	staticK := *merLen

	// create a channel to read from
	reads := make(chan string)

	wg := sync.WaitGroup{}
	// create a new MerMatch
	for i := 0; i+staticK < staticLen; i++ {
		wg.Add(1)
		fmt.Printf("Creating new MerMatch for offset %d (to %d)", i, i+staticK)
		mm := CreateMerMatch(staticK, i, &matchedMers)
		go func() {
			mm.FindMers(reads)
			wg.Done()
		}()
	}

	// read from the file and send to the channel
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txt := scanner.Text()

		reads <- txt

	}

	// close the channel
	close(reads)

	// wait for the goroutines to finish
	wg.Wait()

	// summarise the results
	matchedMers.Summarise()

	// close the file
	f.Close()

}
