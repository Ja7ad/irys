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

	c, err := irys.New(irys.DefaultNode2, matic)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.GetMetaData(context.Background(), "XjzDyneweD_Dmhuaipbi7HyXXvsY6IkMcIsumlB0G2M")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)
}
