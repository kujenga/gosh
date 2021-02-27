// Command stamp provides utilities for easy creation of timestamps.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var stampType string

func main() {
	flag.StringVar(&stampType, "t", "numeric",
		"Type of stamp to create: {numeric|rfc3339}")
	flag.Parse()

	n := time.Now()

	switch stampType {
	case "numeric":
		fmt.Println(n.Format("20060102150405"))
	case "rfc3339":
		fmt.Println(n.Format(time.RFC3339))
	default:
		flag.Usage()
		os.Exit(1)
	}
}
