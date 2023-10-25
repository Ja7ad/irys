package main

import (
	"fmt"
	"github.com/Ja7ad/irys"
	"io"
	"log"
)

func main() {
	matic, err := currency.NewMatic("foo", "bar")
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
	defer file.Data.Close()

	b, err := io.ReadAll(file.Data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b), file.Header, file.ContentLength, file.ContentType)
}
