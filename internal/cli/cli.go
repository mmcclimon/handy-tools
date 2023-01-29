package cli

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func NewFlagSet(name string) *pflag.FlagSet {
	return pflag.NewFlagSet(name, pflag.ContinueOnError)
}

func ParseFlags(fs *pflag.FlagSet) {
	err := fs.Parse(os.Args[1:])

	switch err {
	case nil:
		return
	case pflag.ErrHelp:
		os.Exit(0)
	default:
		fmt.Println(err)
		os.Exit(1)
	}
}
