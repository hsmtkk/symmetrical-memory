package work

import (
	"log"
	"path"
	"strings"

	"github.com/hsmtkk/symmetrical-memory/convert"
)

type Worker interface {
	Run(id int)
}

func New(inDir, outDir string, inFiles chan string, converter convert.Converter) Worker {
	return &workerImpl{inDir, outDir, inFiles, converter}
}

type workerImpl struct {
	inDir     string
	outDir    string
	inFiles   chan string
	converter convert.Converter
}

func (w *workerImpl) Run(id int) {
	log.Printf("start %d worker", id)
	for file := range w.inFiles {
		log.Printf("converting %s file", file)
		inFile := path.Join(w.inDir, file)
		outName := strings.ReplaceAll(file, ".webp", ".jpg")
		outFile := path.Join(w.outDir, outName)
		if err := w.converter.Convert(inFile, outFile); err != nil {
			log.Printf("failed to convert %s file; %s", file, err)
		}
	}
	log.Printf("finish %d worker", id)
}
