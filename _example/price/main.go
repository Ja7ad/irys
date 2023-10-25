package main

import (
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/token"
	"log"
)

func main() {
	matic, err := token.NewMatic("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
	c, err := irys.New(irys.DefaultNode1, matic)
	if err != nil {
		log.Fatal(err)
	}

	p, err := c.GetPrice(100000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p.Int64())
}
