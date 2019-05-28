package root

import (
	flag "github.com/spf13/pflag"
)

func FlagE(f *flag.FlagSet, e *string) {
	f.StringVarP(e, "environment", "e", "", "Environment")
}
