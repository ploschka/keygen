package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/ploschka/keygen/internal/keygen"
)

var (
	fs         *flag.FlagSet = flag.NewFlagSet("keygen", flag.ExitOnError)
	bitFlag                  = fs.Uint64("b", 32, "Number of bytes for the key")
	stringFlag               = fs.String("s", "", "Specify a string to be encoded")
	exportFlag               = fs.Bool("e", false, `add "export" to variable definition`)
	outputFlag               = fs.String("o", "", "Specify file to write keys")
	appendFlag               = fs.Bool("a", false, "Append to file")
)

func errCouldNotParseInt(arg string) error {
	return fmt.Errorf("could not parse int %s", arg)
}

func errCouldNotOpenFile(arg string) error {
	return fmt.Errorf("could not open file %s", arg)
}

func errCouldNotGenerate(arg uint64) error {
	return fmt.Errorf("could not generate key of length %v bytes", arg)
}

func exit(err error) {
	fmt.Fprintln(fs.Output(), err)
}

func isDigits(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func main() {
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s\n", fs.Name())
		fs.PrintDefaults()
	}

	fs.Parse(os.Args[1:])

	var writeTo *os.File

	if len(*outputFlag) > 0 {
		var flag int
		if *appendFlag {
			flag = os.O_WRONLY | os.O_APPEND | os.O_CREATE
		} else {
			flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		}
		var err error
		writeTo, err = os.OpenFile(*outputFlag, flag, 0666)
		if err != nil {
			exit(errors.Join(errCouldNotOpenFile(*outputFlag), err))
		}
		defer writeTo.Close()
	} else {
		writeTo = os.Stdout
	}

	var pattern string
	if *exportFlag {
		pattern = "export %s=%q\n"
	} else {
		pattern = "%s=%q\n"
	}

	count := fs.NArg()

	if count >= 2 {
		if count%2 != 0 {
			count--
		}

		i := 0

		for i < count {
			name := fs.Arg(i)
			arg := fs.Arg(i + 1)

			if isDigits(arg) {
				numarg, err := strconv.ParseUint(arg, 10, 64)
				if err != nil {
					exit(errors.Join(errCouldNotParseInt(arg), err))
				}

				rand, err := keygen.GenerateRand(numarg)
				if err != nil {
					exit(errors.Join(errCouldNotGenerate(numarg), err))
				}

				str := base64.StdEncoding.EncodeToString(rand)
				fmt.Fprintf(writeTo, pattern, name, str)
			} else {
				str := base64.StdEncoding.EncodeToString([]byte(arg))
				fmt.Fprintf(writeTo, pattern, name, str)
			}
			i += 2
		}
	} else {
		if len(*stringFlag) > 0 {
			fmt.Fprintln(writeTo, base64.StdEncoding.EncodeToString([]byte(*stringFlag)))
			return
		}

		rand, err := keygen.GenerateRand(*bitFlag)
		if err != nil {
			exit(errors.Join(errCouldNotGenerate(*bitFlag), err))
		}

		fmt.Fprintln(writeTo, base64.StdEncoding.EncodeToString(rand))
	}
}
