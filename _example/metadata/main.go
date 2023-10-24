package main

import (
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/token"
	"log"
)

func main() {
	matic, err := token.NewMatic("foo")
	if err != nil {
		log.Fatal(err)
	}
	c, err := irys.New(irys.DefaultNode2, matic)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.GetMetaData("XjzDyneweD_Dmhuaipbi7HyXXvsY6IkMcIsumlB0G2M")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)
}
