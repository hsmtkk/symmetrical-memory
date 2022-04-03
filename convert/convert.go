package convert

import (
	"fmt"
	"image/jpeg"
	"os"

	"golang.org/x/image/webp"
)

type Converter interface {
	Convert(inFile, outFile string) error
}

func New() Converter {
	return &converterImpl{}
}

type converterImpl struct{}

func (c *converterImpl) Convert(inFile, outFile string) error {
	inf, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("failed to open %s file to read; %w", inFile, err)
	}
	defer inf.Close()

	img, err := webp.Decode(inf)
	if err != nil {
		return fmt.Errorf("failed to decode %s file; %w", inFile, err)
	}

	outf, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s file to write; %w", outFile, err)
	}
	defer outf.Close()

	if err := jpeg.Encode(outf, img, &jpeg.Options{Quality: 100}); err != nil {
		return fmt.Errorf("failed to encode %s file; %w", outFile, err)
	}

	return nil
}
