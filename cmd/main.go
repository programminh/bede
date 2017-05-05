package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/programminh/bede"
)

var source = flag.String("source", "/var/mnt/fserver/mirax/processor/", "Source folder to parse")
var client = flag.String("client", "client.xml", "Path to client dictionary")
var server = flag.String("server", "server.xml", "Path to source dictionary")

func main() {
	flag.Parse()

	fmt.Println("Generating dictionary of", *source, "to", *client, "and", *server)

	if err := bede.GenDict(*source, *client, *server); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
