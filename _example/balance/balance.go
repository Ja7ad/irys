package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"log"
	"math/big"
)

func main() {
	matic, err := currency.NewMatic(
		"foo",
		"https://example.com")
	if err != nil {
		log.Fatal(err)
	}

	c, err := irys.New(irys.DefaultNode1, matic)
	if err != nil {
		log.Fatal(err)
	}

	b, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.String())

	conf, err := c.TopUpBalance(context.Background(), big.NewInt(2000000000000000))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf)
}
