package main

import (
	"fmt"
	"log"
	"os"

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

	file, err := os.Open("image.jpg")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.Upload(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)
}
