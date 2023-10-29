package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"log"
)

func main() {
	matic, err := currency.NewMatic(
		"foobar",
		"https://exampleRPC.com")
	if err != nil {
		log.Fatal(err)
	}
	c, err := irys.New(irys.DefaultNode1, matic)
	if err != nil {
		log.Fatal(err)
	}

	p, err := c.GetPrice(context.Background(), 100000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p.Int64())
}
