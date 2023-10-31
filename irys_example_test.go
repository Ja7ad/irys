package irys

import (
	"context"
	"fmt"
	"log"

	"github.com/Ja7ad/irys/currency"
)

func ExampleNew() {
	matic, err := currency.NewMatic("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
	c, err := New(DefaultNode1, matic, false)
	if err != nil {
		log.Fatal(err)
	}

	p, err := c.GetPrice(context.Background(), 100000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p.Int64())
}
