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
	c, err := irys.New(irys.DefaultNode1, matic, false)
	if err != nil {
		log.Fatal(err)
	}

	p, err := c.GetPrice(context.Background(), 100000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p.Int64())
}
