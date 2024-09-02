package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/ploschka/keygen/internal/keygen"
)

var (
	bitFlag    = flag.Uint64("b", 32, "Number of bytes for the key")
	stringFlag = flag.String("s", "", "Specify a string to be encoded")
	exportFlag = flag.Bool("e", false, `add "export" to variable definition`)
)

func isDigits(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage of keygen")
		flag.PrintDefaults()
	}

	flag.Parse()

	var pattern string
	if *exportFlag {
		pattern = "export %s=%s\n"
	} else {
		pattern = "%s=%s\n"
	}

	count := flag.NArg()

	if count >= 2 {
		if count%2 != 0 {
			count--
		}

		i := 0

		for i < count {
			name := flag.Arg(i)
			arg := flag.Arg(i + 1)

			if isDigits(arg) {
				numarg, err := strconv.ParseUint(arg, 10, 64)
				if err != nil {
					os.Exit(-1)
				}

				rand, err := keygen.GenerateRand(numarg)
				if err != nil {
					os.Exit(-1)
				}

				str := base64.StdEncoding.EncodeToString(rand)
				fmt.Printf(pattern, name, str)
			} else {
				str := base64.StdEncoding.EncodeToString([]byte(arg))
				fmt.Printf(pattern, name, str)
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
