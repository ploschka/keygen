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
	exportFlag = flag.Bool("e", false, `add "export" to variable definition`)
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage of keygen")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() > 0 {
		for _, v := range flag.Args() {
			rand, err := keygen.GenerateRand(*bitFlag)
			if err != nil {
				os.Exit(-1)
			}

			str := base64.StdEncoding.EncodeToString(rand)

			if *exportFlag {
				fmt.Printf("export %s=%s\n", v, str)
			} else {
				fmt.Printf("%s=%s\n", v, str)
			}
		}
	} else {
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
}
