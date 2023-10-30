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
		"foobar",
		"https://exampleRPC.com")
	if err != nil {
		log.Fatal(err)
	}

	c, err := irys.New(irys.DefaultNode1, matic, true)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("/home/javad/Downloads/Shotgun_Types.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.ChunkUpload(context.Background(), file, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)

}
