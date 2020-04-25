package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	host = "golang.org"
)

func main() {
	dial()
	lookupHost()
	lookupAddr()
}

func dial() {
	defer timeTrack(time.Now(), "dial")
	conn, err := net.Dial("tcp", host+":80")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
}

func lookupHost() {
	defer timeTrack(time.Now(), "lookupHost")
	// LookupHost looks up the given host using the local resolver.
	// It returns a slice of that host's addresses.
	_, err := net.LookupHost(host)
	if err != nil {
		log.Printf("failed to lookup IP of domain %s.", host)
		os.Exit(1)
	}
}

func lookupAddr() {
	defer timeTrack(time.Now(), "lookupAddr")
	// LookupAddr performs a reverse lookup for the given address,
	// returning a list of names mapping to that address.
	_, err := net.LookupAddr("6.8.8.8")
	if err != nil {
		log.Println(err)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
