# irys-go [![Go Reference](https://pkg.go.dev/badge/github.com/Ja7ad/irys.svg)](https://pkg.go.dev/github.com/Ja7ad/irys)
Go Implementation SDK of Irys network, irys is the only provenance layer. It enables users to scale permanent data and precisely attribute its origin (arweave bundlr).

| Currency          | arweave | ethereum | matic | bnb | avalanche | solana | arbitrum | fantom | near | algorand | aptos |
|-------------------|---------|----------|-------|-----|-----------|--------|----------|--------|------|----------|-------|
| Price API         |    x    |     x    |   x   |  x  |     x     |    -   |     x    |    x   |   -  |     -    |   -   |
| Balance API       |    x    |     x    |   x   |  x  |     x     |    -   |     x    |    x   |   -  |     -    |   -   |
| Upload File API   |    -    |     x    |   x   |  x  |     x     |    -   |     x    |    x   |   -  |     -    |   -   |
| Chunk File API    |    -    |     x    |   x   |  x  |     x     |    -   |     x    |    x   |   -  |     -    |   -   |
| Upload Folder API |    -    |     -    |   -   |  -  |     -     |    -   |     -    |    -   |   -  |     -    |   -   |
| Widthdraw API     |    -    |     -    |   -   |  -  |     -     |    -   |     -    |    -   |   -  |     -    |   -   |
| Get Receipt API   |    -    |     x    |   x   |  x  |     x     |    -   |     x    |    x   |   -  |     -    |   -   |
| Verify Receipt API |    -    |     -    |   -   |  -  |     -     |    -   |     -    |    -   |   -  |     -    |   -   |
| Found API         |    -    |     x    |   x   |  x  |     x     |    -   |     x    |    x   |   -  |     -    |   -   |

## Install

```shell
go get -u  github.com/Ja7ad/irys
```

## Examples

[example](_example) of irys sdk in golang 

### Upload

```go
package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"log"
	"os"
)

func main() {
	matic, err := currency.NewMatic(
		"foo",
		"https://exampleRPC.com")
	if err != nil {
		log.Fatal(err)
	}

	c, err := irys.New(irys.DefaultNode2, matic, false)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("image.jpeg")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.BasicUpload(context.Background(), file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)

}
```

### Download

```go
package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"io"
	"log"
)

func main() {
	matic, err := currency.NewMatic("foo", "https://exampleRPC.com")
	if err != nil {
		log.Fatal(err)
	}
	c, err := irys.New(irys.DefaultNode2, matic, false)
	if err != nil {
		log.Fatal(err)
	}

	file, err := c.Download(context.Background(), "XjzDyneweD_Dmhuaipbi7HyXXvsY6IkMcIsumlB0G2M")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Data.Close()

	b, err := io.ReadAll(file.Data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b), file.Header, file.ContentLength, file.ContentType)
}
```

### Calculate Price

```go
package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"log"
)

func main() {
	matic, err := currency.NewMatic("foo", "https://exampleRPC.com")
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
```

### Get Metadata

```go
package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/currency"
	"log"
)

func main() {
	matic, err := currency.NewMatic("foo", "https://exampleRPC.com")
	if err != nil {
		log.Fatal(err)
	}
	c, err := irys.New(irys.DefaultNode2, matic, false)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := c.GetMetaData(context.Background(), "XjzDyneweD_Dmhuaipbi7HyXXvsY6IkMcIsumlB0G2M")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx)
}
```

## Todo

- [x] arweave network
- [x] ethereum network
- [x] polygon matic network
- [x] concurrent and chunk upload
- [x] get chunk upload transaction response
- [ ] fix bug finish chunk upload for finalizing
- [ ] unit test
- [x] found API
- [ ] upload folder
- [ ] withdraw balance
- [x] get loaded balance