# irys-go [![Go Reference](https://pkg.go.dev/badge/github.com/Ja7ad/irys.svg)](https://pkg.go.dev/github.com/Ja7ad/irys)
Go Implementation SDK of Irys network, irys is the only provenance layer. It enables users to scale permanent data and precisely attribute its origin (arweave bundlr).

## Install

```shell
go get -u  github.com/Ja7ad/irys
```

## Examples

example of irys sdk in golang 

### Upload

```go
package main

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/token"
	"log"
	"os"
)

func main() {
	matic, err := token.NewMatic(
		"foo",
		"bar")
	if err != nil {
		log.Fatal(err)
	}

	c, err := irys.New(irys.DefaultNode2, matic)
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
	"fmt"
	"github.com/Ja7ad/irys"
	"github.com/Ja7ad/irys/token"
	"io"
	"log"
)

func main() {
	matic, err := token.NewMatic("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
	c, err := irys.New(irys.DefaultNode2, matic)
	if err != nil {
		log.Fatal(err)
	}

	file, err := c.Download("XjzDyneweD_Dmhuaipbi7HyXXvsY6IkMcIsumlB0G2M")
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
```

### Get Metadata

```go
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
```

## Todo

- [x] arweave token
- [x] ethereum token
- [x] polygon matic token
- [x] fantom token
- [x] bnb token
- [x] arbitrum token
- [x] avalanche token
- [ ] concurrent and chunk upload
- [ ] unit test
- [x] found API
- [ ] upload folder
- [ ] withdraw balance
- [x] get loaded balance