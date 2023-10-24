package imageHandler

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"
)

type FileType string

const (
	JPG FileType = "jpg"
	PNG FileType = "png"
)

func ServeFrames(imgByte []byte, address string, fileType FileType) {
	switch fileType {
	case JPG:
		img, _, err := image.Decode(bytes.NewReader(imgByte))
		if err != nil {
			log.Fatalln(err)
		}

		out, _ := os.Create(address)
		defer func(out *os.File) {
			err := out.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(out)

		var opts jpeg.Options
		opts.Quality = 1

		err = jpeg.Encode(out, img, &opts)
		if err != nil {
			log.Println(err)
		}
	}
}
