# go-swiss-qr-bill

A small package to parse Swiss QR-Bills, according to the 
[Swiss Payment Standards specification](https://www.six-group.com/dam/download/banking-services/standardization/qr-bill/ig-qr-bill-v2.2-en.pdf).


## Usage

### CLI

```bash
go install github.com/denysvitali/go-swiss-qr-bill/cmd/swiss-qr-bill@master
swiss-qr-bill < ./samples/1.txt
```

### Package

```go
package main
import (
	swissqrcode "github.com/denysvitali/go-swiss-qr-bill"
)


func main(){
    qrCode, err := swissqrcode.Decode("...")
    if err != nil {
        // ...
    }	
	// ...
}
```