package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/configs"
	"github.com/Ja7ad/irys/currency"
	"log"
	"os"
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

	file, err := os.Open("absolute_path_to_file")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.ChunkUpload(context.Background(), file, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)

}
