package main

import (
	"log"
	"os"
	"sync"

	"github.com/hsmtkk/symmetrical-memory/convert"
	"github.com/hsmtkk/symmetrical-memory/work"
	"github.com/spf13/cobra"
)

const defaultNumGoRoutines = 4

var (
	numGoRoutines int
)

func main() {
	command := &cobra.Command{
		Use:  "symmetrical-memory inDir outDir",
		Run:  run,
		Args: cobra.ExactArgs(2),
	}
	command.Flags().IntVar(&numGoRoutines, "numGoRoutines", defaultNumGoRoutines, "number of go routines")
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	inDir := args[0]
	outDir := args[1]
	inFiles := make(chan string)
	worker := work.New(inDir, outDir, inFiles, convert.New())
	var wg sync.WaitGroup
	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			worker.Run(i)
		}(i)
	}
	entries, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		inFiles <- entry.Name()
	}
	close(inFiles)
	wg.Wait()
}
