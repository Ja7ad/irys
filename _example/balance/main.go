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

	c, err := irys.New(irys.DefaultNode1, matic, true)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	b, err := c.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.String())
}
