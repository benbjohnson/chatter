package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/benbjohnson/chatter"
)

func main() {
	// Parse command line arguments.
	addr := flag.String("addr", ":9000", "bind address")
	flag.Parse()

	// Create and listen.
	h := &chatter.Handler{}
	fmt.Printf("Listening on http://localhost%s\n", *addr)
	http.ListenAndServe(*addr, h)
}
