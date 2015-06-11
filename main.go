package main

import (
	"github.com/hni/gsafeeder/lib"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) != 3 {
       fmt.Fprintf(os.Stderr, "Usage: %s <gsa> <xml feed file>\n", os.Args[0])
       os.Exit(1)
    }
    gsafeeder.Upload(os.Args[1], os.Args[2])
}
