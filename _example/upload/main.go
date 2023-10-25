package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"log"
	"os"
)

func main() {
	matic, err := currency.NewMatic(
		"foo",
		"bar")
	if err != nil {
		log.Fatal(err)
	}

	c, err := irys.New(irys.DefaultNode2, matic)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("image.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.BasicUpload(context.Background(), file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)

}
