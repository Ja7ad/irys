package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"log"
	"os"
	"time"
)

func main() {
	matic, err := currency.NewMatic(
		"foobar",
		"https://exampleRPC.com")
	if err != nil {
		log.Fatal(err)
	}

	c, err := irys.New(irys.DefaultNode1, matic, false)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("/home/javad/Downloads/mtn.json")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := c.BasicUpload(ctx, file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)

}
