package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"
)

func establishTCPConnection(address string) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Printf("Failed to resolve TCP address: %v", err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v", err)
		return conn, err
	}
	return conn, nil
}

func readResponse(r *bufio.Reader) error {
	httpStatus, err := r.ReadString('\n')
	if err != nil {
		if errors.Is(err, syscall.ECONNRESET) {
			fmt.Printf("Connection was reset by the server: %v \n", err)
			return err
		}

		if errors.Is(err, io.EOF) {
			fmt.Printf("Nothing to read: %v \n", err)
			return err
		}

		fmt.Printf("Unexpected error: %v", err)
		return err
	}
	fmt.Print(httpStatus)
	bytesRemaining := r.Buffered()
	fmt.Printf("Bytes available: %d\n", bytesRemaining)
	return nil
}

func isAliveAfter(address string, d time.Duration) bool {
	conn, err := establishTCPConnection(address)
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\nConnection: keep-alive\r\n\r\n")
	err = readResponse(reader)
	if err != nil {
		fmt.Printf("Reading response to initial request failed: %v", err)
		return false
	}

	fmt.Printf("Waiting %ds\n", int(d.Seconds()))
	time.Sleep(d)

	if err != nil {
		fmt.Printf("Reading response to second request failed: %v", err)
		return false
	}

	return true
}

func main() {

	address := flag.String("address", "www.google.com:80", "Address of test target in form of host:port. Example: google.com:80 or localhost:80")
	initialWait := flag.Uint("period", 10, "Duration in seconds to wait before sending second request.")
	interval := flag.Int("interval", 10, "Number of seconds to add to the wait time each test run.")
	flag.Parse()

	n := time.Duration(*initialWait) * time.Second

	for isAliveAfter(*address, n) {
		if n > 350*time.Second {
			break
		}
		n += time.Duration(*interval) * time.Second

	}

	fmt.Printf("Keepalive for %s is shorter than %d seconds.\n", *address, int(n.Seconds()))
}
