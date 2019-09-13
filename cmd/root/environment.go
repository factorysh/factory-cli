package root

import (
	flag "github.com/spf13/pflag"
)

func FlagE(f *flag.FlagSet) {
	f.StringVarP(&Environment, "environment", "e", "", "Environment name to target")
}
