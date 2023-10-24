package main

import (
	"fmt"
	"io"
	"log"

	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/token"
)

func main() {
	matic, err := token.NewMatic("foo")
	if err != nil {
		log.Fatal(err)
	}
	c, err := irys.New(irys.DefaultNode2, matic)
	if err != nil {
		log.Fatal(err)
	}

	file, err := c.Download("XjzDyneweD_Dmhuaipbi7HyXXvsY6IkMcIsumlB0G2M")
	if err != nil {
		log.Fatal(err)
	}
	defer func(Data io.ReadCloser) {
		err := Data.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file.Data)

	b, err := io.ReadAll(file.Data)
	if err != nil {
		log.Fatal(err)
	}
	// imageHandler.ServeFrames(b, "./_example/download/output.jpeg", imageHandler.JPG)
	fmt.Println(string(b), file.Header, file.ContentLength, file.ContentType)
}
