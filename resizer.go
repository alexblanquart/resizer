package main

import (
	"flag"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"strconv"
)

func resizer(target string, width uint, height uint) filepath.WalkFunc {
	return func(img string, f os.FileInfo, err error) error {
		if _, err := os.Stat(img); os.IsExist(err) {
			os.Remove(img)
			log.Printf("%s image deleted to be renewed", img)
		}

		file, err := os.Open(img)
		defer file.Close()

		initial, _, err := image.Decode(file)
		if err != nil {
			// Happen for root directory and images with unkown/unhandled format.
			// For now : png and jpeg!
			return nil
		}

		// resize to specified width and height using Lanczos resampling
		// and preserve aspect ratio
		resized := resize.Resize(width, height, initial, resize.Lanczos3)

		// get name of image without the extra stuff to compute new path
		name := filepath.Base(img)
		ext := filepath.Ext(img)
		nameWithoutExt := name[:len(name)-len(ext)]
		newPath := target + nameWithoutExt + ".jpg"

		// create the new file
		out, err := os.Create(newPath)
		if err != nil {
			log.Printf("%s when creating new file", err)
			return err
		}
		defer out.Close()

		// write new image to file
		return jpeg.Encode(out, resized, nil)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 4 {
		log.Printf("usage : resizer directoryWithOriginalImages directoryWithResizedImages widthForResizing heightForResizing")
		os.Exit(2)
	}
	original := flag.Arg(0)
	target := flag.Arg(1)
	width, _ := strconv.Atoi(flag.Arg(2))
	height, _ := strconv.Atoi(flag.Arg(3))

	log.Printf("Resizing images inside %s into %s with width %d and height %d", original, target, width, height)

	// create directory if not exists
	if _, err := os.Stat(target); os.IsNotExist(err) {
		err := os.Mkdir(target, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	filepath.Walk(original, resizer(target, uint(width), uint(height)))
}
