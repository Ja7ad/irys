package irys

import (
	"fmt"
	"github.com/Ja7ad/irys/token"
	"log"
)

func ExampleNew() {
	matic, err := token.NewMatic("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
	c, err := New(DefaultNode1, matic)
	if err != nil {
		log.Fatal(err)
	}

	p, err := c.GetPrice(100000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(p.Int64())
}
