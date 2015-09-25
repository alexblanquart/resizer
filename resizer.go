package main

import (
	"flag"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/davecheney/profile"
	"github.com/nfnt/resize"
	"strconv"
	"sync"
)

// Resize the JPEG image denoted by its path with specified width and height.
// As stipulated inside resize library, one of the two can be omitted.
// The newly JPEG resized image is in the end in the target directory with the same name.
// In case the resized image is already present in the target directory, it is renewed by removing the old one before.
func resizeImage(img string, target string, width uint, height uint) error {
	// get name of image without the extra stuff to compute new path
	name := filepath.Base(img)
	ext := filepath.Ext(img)
	nameWithoutExt := name[:len(name)-len(ext)]
	newPath := target + nameWithoutExt + ".jpg"

	// verify image is not yet present
	if _, err := os.Stat(newPath); os.IsExist(err) {
		os.Remove(newPath)
		log.Printf("%s image deleted to be renewed", newPath)
	}

	file, err := os.Open(img)
	defer file.Close()
	initial, _, err := image.Decode(file)
	if err != nil {
		// Happen for root directory and images with unkown/unhandled format
		// Only JPEG format for now
		return nil
	}

	// resize to specified width and height using Lanczos resampling
	// and preserve aspect ratio
	resized := resize.Resize(width, height, initial, resize.Lanczos3)

	// create the new file
	out, err := os.Create(newPath)
	defer out.Close()
	if err != nil {
		log.Printf("%s when creating new file", err)
		return err
	}

	// write new image to file
	return jpeg.Encode(out, resized, nil)
}

func main() {
	defer profile.Start(&profile.Config{
        MemProfile: true,
        ProfilePath: ".",
    }).Stop()

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

	// create directory if not yet exist
	if _, err := os.Stat(target); os.IsNotExist(err) {
		err := os.Mkdir(target, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	// resize images
	var wg sync.WaitGroup
	filepath.Walk(original, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			resizeImage(path, target, uint(width), uint(height))
			log.Printf("%s image resized", path)
		}()
		return nil
	})

	wg.Wait()
}
