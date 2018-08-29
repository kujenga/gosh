// Command urlencode provides a utliity for url encoding content.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error reading from stdin: %v", err)
	}

	out := url.PathEscape(string(data))
	fmt.Fprint(os.Stdout, out)
}
