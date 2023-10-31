package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

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

	err = c.TopUpBalance(ctx, big.NewInt(321000000000023))
	if err != nil {
		log.Fatal(err)
	}

	balance, err := c.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance)
}
