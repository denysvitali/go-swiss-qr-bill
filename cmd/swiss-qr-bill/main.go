package main

import (
	"encoding/json"
	"fmt"
	swissqrcode "github.com/denysvitali/go-swiss-qr-bill"
	"io"
	"os"
)

func main() {
	qrCodeText, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read all from stdin: %v", err)
		os.Exit(1)
	}

	qrCode, err := swissqrcode.Decode(string(qrCodeText))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode QR Code: %v", err)
		os.Exit(2)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(qrCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write JSON to stdout: %v", err)
		os.Exit(3)
	}
}
