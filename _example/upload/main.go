package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
)

func main() {
	matic, err := currency.NewMatic("ExamplePrivateKey", "ExampleRpc")
	if err != nil {
		log.Fatal(err)
	}

	c, err := irys.New(irys.DefaultNode1, matic, false)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("absolute_path_to_file")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	tx, err := c.BasicUpload(ctx, bs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)
}
