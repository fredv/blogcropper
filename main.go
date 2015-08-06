package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/oliamb/cutter"
)

const SOI_MARKER = 0xd8
const (
	imagePath = "static/images/posts"
	imageResX = 1200
	imageResY = 630
)

type imageTarget struct {
	name string
	resX int
	resY int
}

var targets []imageTarget = []imageTarget{
	{"facebook_square_small", 200, 200},
	{"facebook_square", 630, 630},
	{"facebook_share", 1200, 630},
	{"blog_large", 617, 242},
	{"blog_thumb", 293, 118},
	{"twitter_large", 1200, 626},
}

func maxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func createBaseImage(img image.Image) (image.Image, error) {
	var err error
	bounds := img.Bounds()
	dx, dy := bounds.Dx(), bounds.Dy()
	switch {
	case dx == imageResX && dy == imageResY:
		return img, nil // nothing needs to be done
	default:
		target := img
		// resize accordingly
		switch {
		case dx < imageResX:
			dy = dy * imageResX / dx
			target = imaging.Resize(target, imageResX, dy, imaging.Lanczos)
		case dx > imageResX:
			dy = maxInt(dy, imageResY)
			target, err = cutter.Crop(target, cutter.Config{
				Width:   imageResX,
				Height:  dy,
				Mode:    cutter.Centered,
				Options: cutter.Copy,
			})
			if err != nil {
				return nil, err
			}
		}

		switch {
		case dy < imageResY:
			// add blurred image to top and bottom
			largeImg := imaging.Resize(target, imageResX, imageResY, imaging.Lanczos)
			backgroundImg := imaging.Blur(largeImg, 50)
			target = imaging.Overlay(backgroundImg, target, image.Pt(0, (imageResY-dy)/2), 1.0)
		case dy > imageResY:
			target, err = cutter.Crop(target, cutter.Config{
				Width:   imageResX,
				Height:  imageResY,
				Mode:    cutter.Centered,
				Options: cutter.Copy,
			})
			if err != nil {
				return nil, err
			}
		}

		return target, nil
	}
}

func resizeImage(filename string, img image.Image, resX, resY int) error {
	fmt.Println(filename)
	var err error
	res := img

	if resX < imageResX || resY < imageResY {
		ratioX, ratioY := float64(resX)/float64(imageResX), float64(resY)/float64(imageResY)
		if ratioX < ratioY {
			img = imaging.Resize(img, int(imageResX*ratioY), resY, imaging.Lanczos)
		} else {
			img = imaging.Resize(img, resX, int(imageResY*ratioX), imaging.Lanczos)
		}
		res, err = cutter.Crop(img, cutter.Config{
			Width:   resX,
			Height:  resY,
			Mode:    cutter.Centered,
			Options: cutter.Copy,
		})
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	err = png.Encode(w, res)
	if err != nil {
		return err
	}
	return w.Flush()
}

func handleFile(fileName string) error {
	var baseImage image.Image

	for i := range targets {
		filename := strings.Replace(fileName, imagePath+"/", imagePath+"/gen/"+targets[i].name+"_", 1)
		_, err := os.Stat(filename)
		switch {
		case err == nil:
			// nothing should be done
		case os.IsNotExist(err):
			if baseImage == nil {
				r, err := os.Open(fileName)
				if err != nil {
					return err
				}

				img, err := png.Decode(r)
				r.Close()
				if err != nil {
					return err
				}

				baseImage, err = createBaseImage(img)
				if err != nil {
					return err
				}
			}
			err = resizeImage(filename, baseImage, targets[i].resX, targets[i].resY)
			if err != nil {
				return err
			}
		default:
			return err
		}
	}
	return nil
}

func main() {
	pattern := imagePath + "/*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(imagePath + "/gen"); os.IsNotExist(err) {
		err = os.Mkdir(imagePath+"/gen", 0775)
		if err != nil {
			log.Fatal(err)
		}
	}

	wg := new(sync.WaitGroup)
	fileChan := make(chan string)
	errChan := make(chan error)

	foundError := false
	go func(e <-chan error) {
		for err := range e {
			fmt.Fprintf(os.Stderr, "ERROR: %s", err)
			foundError = true
		}
	}(errChan)

	for i := runtime.NumCPU(); i > 0; i-- {
		wg.Add(1)
		go func(c <-chan string, e chan<- error, wg *sync.WaitGroup) {
			defer wg.Done()
			for filename := range c {
				err := handleFile(filename)
				if err != nil {
					e <- err
				}
			}
		}(fileChan, errChan, wg)
	}

	for _, fileName := range matches {
		if !strings.Contains(fileName, ".png") {
			continue
		}

		fileChan <- fileName
	}

	close(fileChan)
	wg.Wait()
	close(errChan)

	if foundError {
		os.Exit(1)
	}
}
