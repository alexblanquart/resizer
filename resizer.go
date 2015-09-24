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

func resizer(original string, f os.FileInfo, err error) error {
	file, err := os.Open(original)
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
	name := filepath.Base(original)
	ext := filepath.Ext(original)
	nameWithoutExt := name[:len(name)-len(ext)]
	newPath := target + nameWithoutExt + ".jpg"
	log.Printf("%s image about to created", newPath)

	// create the new file
	out, err := os.Create(newPath)
	if err != nil {
		log.Printf("%s when creating new file", err)
		return err
	}
	defer out.Close()

	// write new image to file
	return jpeg.Encode(out, resized, nil)
	if err != nil {
		log.Printf("%s when encoding new file as jpeg image", err)
	}
	return nil
}

func resize(original string, target string, width uint, height uint){
	// call resizer on each walked file
	filepath.Walk(original, resizer)
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

	// renew resize directory no matter what
	os.RemoveAll(target)
	err := os.Mkdir(target, 0777)
	if err != nil {
		log.Fatal(err)
	}

	resize(original, target, width, height)
}
