package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/configs"
	"github.com/Ja7ad/irys/currency"
)

func main() {
	matic, err := currency.NewMatic(configs.ExamplePrivateKey, configs.ExampleRpc)
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
