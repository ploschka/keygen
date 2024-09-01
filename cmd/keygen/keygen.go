package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/ploschka/keygen/internal/keygen"
)

var (
	bitFlag    = flag.Uint64("b", 32, "Number of bytes for the key")
	stringFlag = flag.String("s", "", "Specify a string to be encoded")
)

func main() {
	flag.Parse()

	if len(*stringFlag) > 0 {
		fmt.Println(base64.StdEncoding.EncodeToString([]byte(*stringFlag)))
		return
	}

	rand, err := keygen.GenerateRand(*bitFlag)
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(rand))
}
