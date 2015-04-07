package main

import (
	"bufio"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/oliamb/cutter"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

const SOI_MARKER = 0xd8

func main() {
	pattern := "./static/images/posts/*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		panic(err.Error())
	}
	for _, fileName := range matches {
		if !strings.Contains(fileName, ".png") {
			continue
		}
		fmt.Println(fileName)
		r, err := os.Open("./" + fileName)
		if err != nil {
			panic(err.Error())
		}
		defer r.Close()
		img, err := png.Decode(r)
		if err != nil {
			panic(err.Error())
		}
		croppedImg, err := cutter.Crop(img, cutter.Config{
			Width:   191,
			Height:  191,
			Mode:    cutter.Centered,
			Options: cutter.Copy,
		})

		if err != nil {
			panic(err.Error())
		}

		filename := strings.Replace(fileName, "static/images/posts/", "static/images/posts/gen/facebook_square_small_", 1)
		if _, err = os.Stat(filename); os.IsNotExist(err) {
			squareImg := imaging.Resize(croppedImg, 200, 200, imaging.Lanczos)
			f, err := os.Create(filename)
			defer f.Close()
			w := bufio.NewWriter(f)

			err = png.Encode(w, squareImg)
			if err != nil {
				panic(err.Error())
			}
			w.Flush()
		}

		filename = strings.Replace(fileName, "static/images/posts/", "static/images/posts/gen/facebook_square_", 1)
		if _, err = os.Stat(filename); os.IsNotExist(err) {
			squareImg := imaging.Resize(croppedImg, 628, 628, imaging.Lanczos)
			f, err := os.Create(filename)
			defer f.Close()
			w := bufio.NewWriter(f)

			err = png.Encode(w, squareImg)
			if err != nil {
				panic(err.Error())
			}
			w.Flush()
		}

		filename = strings.Replace(fileName, "static/images/posts/", "static/images/posts/gen/blog_large_", 1)
		if _, err = os.Stat(filename); os.IsNotExist(err) {
			largeImg := imaging.Resize(img, 617, 247, imaging.Lanczos)
			backgroundImg := imaging.Blur(largeImg, 50)

			dstImg := imaging.Overlay(backgroundImg, img, image.Pt(0, 28), 1.0)

			f, err := os.Create(filename)
			defer f.Close()
			w := bufio.NewWriter(f)

			err = png.Encode(w, dstImg)
			if err != nil {
				panic(err.Error())
			}
			w.Flush()
		}

		filename = strings.Replace(fileName, "static/images/posts/", "static/images/posts/gen/blog_thumb_", 1)
		if _, err = os.Stat(filename); os.IsNotExist(err) {
			largeImg := imaging.Resize(img, 617, 247, imaging.Lanczos)
			backgroundImg := imaging.Blur(largeImg, 50)

			dstImg := imaging.Overlay(backgroundImg, img, image.Pt(0, 28), 1.0)
			blogThumb := imaging.Resize(dstImg, 293, 118, imaging.Lanczos)
			f, err := os.Create(filename)
			defer f.Close()
			w := bufio.NewWriter(f)

			err = png.Encode(w, blogThumb)
			if err != nil {
				panic(err.Error())
			}
			w.Flush()
		}

		filename = strings.Replace(fileName, "static/images/posts/", "static/images/posts/gen/twitter_large_", 1)

		if _, err = os.Stat(filename); os.IsNotExist(err) {
			largeImg := imaging.Resize(img, 1200, 643, imaging.Lanczos)
			spriteImg := imaging.Resize(img, 1200, 317, imaging.Lanczos)
			backgroundImg := imaging.Blur(largeImg, 50)

			dstImg := imaging.Overlay(backgroundImg, spriteImg, image.Pt(0, 163), 1.0)

			f, err := os.Create(filename)
			defer f.Close()
			w := bufio.NewWriter(f)

			err = png.Encode(w, dstImg)
			if err != nil {
				panic(err.Error())
			}
			w.Flush()
		}

		filename = strings.Replace(fileName, "static/images/posts/", "static/images/posts/gen/facebook_share_", 1)

		if _, err = os.Stat(filename); os.IsNotExist(err) {
			largeImg := imaging.Resize(img, 1200, 643, imaging.Lanczos)
			spriteImg := imaging.Resize(img, 1200, 317, imaging.Lanczos)
			backgroundImg := imaging.Blur(largeImg, 50)

			dstImg := imaging.Overlay(backgroundImg, spriteImg, image.Pt(0, 163), 1.0)
			croppedImg, err = cutter.Crop(dstImg, cutter.Config{
				Width:   1200,
				Height:  628,
				Mode:    cutter.Centered,
				Options: cutter.Copy,
			})

			if err != nil {
				panic(err.Error())
			}

			f, err := os.Create(filename)
			defer f.Close()
			w := bufio.NewWriter(f)

			err = png.Encode(w, croppedImg)
			if err != nil {
				panic(err.Error())
			}
			w.Flush()
		}
	}
}
